package autowire

import (
	"io"
	"log"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var dependencies map[string]interface{}
var currentProfile = getProfile()
var ch = make(chan os.Signal)

const Tag = "autowire"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Init Autowire Context")
	signal.Notify(ch, os.Interrupt, os.Kill)
	dependencies = make(map[string]interface{})
}

// InitProd executes function only in production so it is preventing execution in go tests.
// This flexibility could help if we want to skip autowiring struct in our tests.
func InitProd(initFunc func()) {
	if currentProfile == _Production {
		initFunc()
	}
}

// Autowire function injects all dependencies for the given structure v.
// In order dependencies to be injected the desired struct fields should be marked with
// autowire tag.
//
// Injection of concrete type
//
// Currently both type of fields exported and unexported are supported.
// Following snippet shows injecting dependencies inside private structure fields using `autowire:""` tag:
//  type Application struct {
//      config  *configuration.ApplicationConfig `autowire:""`
//  }
//  func (a *Application) SetConfig(config  *configuration.ApplicationConfig)  {
//      a.config = config
//  }
// If we need dependency to be injected into unexported field Set<FieldName> function is required, as show above.
//  type Application struct {
//      Config  *configuration.ApplicationConfig `autowire:""`
//  }
// Injection of dependency into exported is supported also and there is no need to provide additional Setter.
//
// Injection of interface
//
// Often it`s better to rely on interface instead of concrete type, in order to accomplish this decoupling we specify
// interfaces as a type of struct fields. The following snippet demonstrate that
//  type UserService struct {
//      userRoleRepository UserRoleRepository `autowire:"service/InMemoryUserRoleRepository"`
//  }
//  func (u *UserService) SetUserRoleRepository(userRoleRepository UserRoleRepository) {
//      u.userRoleRepository = userRoleRepository
//  }
// UserRoleRepository is simply an interface and InMemoryUserRoleRepository is a struct, which implements this interface.
// As it`s done dependency injection on unexported field, providing Setter is required. Just to highlight unexported fields
// needs Setter while exported not. For more information take a look at example package.
func Autowire(v interface{}) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Invalid:
		log.Println("invalid reflection type")
	case reflect.Ptr:
		structType := getStructPtrFullPath(value)
		_, ok := dependencies[structType]
		if ok {
			log.Printf("%s already autowired... ignored", structType)
		} else {
			log.Printf("Autowiring %s", structType)
			autowireDependencies(value)
			dependencies[structType] = v
		}
	default: // reflect.Array, reflect.Struct, reflect.Interface
		log.Println(value.Type().String() + " value")
	}
	//log.Println(dependencies)
}

// Autowired function returns fully initialized with all dependencies instance, which is ready to be used.
//As the result is empty interface, type assertions is required before using the instance.
//Take a look at https://golang.org/ref/spec#Type_assertions for more information.
// The following snippet demonstrate how could be done :
// 	 app := autowire.Autowired(app.Application{}).(*app.Application)
func Autowired(v interface{}) interface{} {
	value := reflect.ValueOf(v)
	var path string
	switch value.Kind() {
	case reflect.Invalid:
		log.Println("invalid")
	case reflect.Struct:
		path = getFullPath(value.Type().PkgPath(), value.Type().String())
	case reflect.Ptr:
		path = getStructPtrFullPath(value)
	default:
		log.Panicln("Unknown Autowired Typed!")
	}
	dependency, ok := dependencies[path]
	if ok {
		return dependency
	}
	return nil
}

func getStructPtrFullPath(value reflect.Value) string {
	return getFullPath(value.Elem().Type().PkgPath(), value.Type().String())
}

func getFullPath(pkgPath string, typePath string) string {
	var re = regexp.MustCompile(`^\*?.+\.`)
	return pkgPath + re.ReplaceAllString(typePath, "/")
}

func autowireDependencies(value reflect.Value) {
	elem := value.Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		tag, ok := field.Tag.Lookup(Tag)
		if ok {
			var t reflect.Value
			if tag != "" {
				currentDep := findDependency(tag)
				if len(currentDep) == 0 {
					msg := "Unknown dependency " + tag + " found none"
					if currentProfile != _Testing {
						log.Panicln(msg)
					} else {
						log.Println(msg + ", ready for spy")
					}
				} else {
					v := reflect.ValueOf(currentDep[0])
					if v.Type().Implements(field.Type) {
						t = reflect.New(v.Type())
						dependency := currentDep[0]
						setValue(value, elem, i, dependency)
					} else {
						log.Panicln(v.Type().String() + " doesnt Implements: " + field.Type.String())
					}
				}
			} else {
				t = reflect.New(elem.Type().Field(i).Type.Elem())
				dependency, found := dependencies[getStructPtrFullPath(t)]
				if found {
					runeName := []rune(elem.Type().Field(i).Name)
					exported := unicode.IsUpper(runeName[0])
					if exported {
						elem.Field(i).Set(reflect.ValueOf(dependency))
					} else {
						runeName[0] = unicode.ToUpper(runeName[0])
						methodName := "Set" + string(runeName)
						method := value.MethodByName(methodName)
						if method.IsValid() {
							method.Call([]reflect.Value{reflect.ValueOf(dependency)})
						}
					}
				}
			}
		}
	}
}

func setValue(value reflect.Value, elem reflect.Value, i int, dependency interface{}) {
	runeName := []rune(elem.Type().Field(i).Name)
	exported := unicode.IsUpper(runeName[0])
	if exported {
		elem.Field(i).Set(reflect.ValueOf(dependency))
	} else {
		runeName[0] = unicode.ToUpper(runeName[0])
		methodName := "Set" + string(runeName)
		method := value.MethodByName(methodName)
		if method.IsValid() {
			method.Call([]reflect.Value{reflect.ValueOf(dependency)})
		}
	}
}

func findDependency(tagDependencyType string) []interface{} {
	var result []interface{}
	for tmp, dep := range dependencies {
		if strings.Contains(tmp, tagDependencyType) {
			result = append(result, dep)
		}
	}
	return result
}

// Run executes function function after the function execution completes for each autowired struct
// which implements io.Closer interface will be invoked Call() function, so currently active
//occupied resources(connections, channels, etc.) could be released.
func Run(function func()) {
	defer close()
	function()
}

func close() {
	log.Println("Closing...")
	for _, dependency := range dependencies {
		valueDepend := reflect.ValueOf(dependency)
		closerType := reflect.TypeOf((*io.Closer)(nil)).Elem()
		if valueDepend.Type().Implements(closerType) {
			err := dependency.(io.Closer).Close()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}
