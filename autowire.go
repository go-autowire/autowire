package autowire

import (
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var (
	dependencies   map[string]interface{}
	currentProfile = getProfile()
)

// Tag name
const Tag = "autowire"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Init Autowire Context")
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
// Autowire function should be executed in the init() function of the
// package. In order dependencies to be injected the desired struct fields
// should be marked with autowire tag.
//
// Concrete type Injection
//
// Currently both type of fields exported and unexported are supported.
//
// Unexported field
//
// Following snippet shows injecting dependencies inside private structure
// fields using `autowire:""` tag:
//  type Application struct {
//      config  *configuration.ApplicationConfig `autowire:""`
//  }
//  func (a *Application) SetConfig(config  *configuration.ApplicationConfig)  {
//      a.config = config
//  }
// If we need given dependency to be injected into unexported field Setter
// function(Set<FieldName>) is required, as show above. If you have a field
// called config (lower case, unexported), the setter method should be called
// SetConfig, not Setconfig.
//
// Exported field
//
// The use of upper-case names of the struct fields indicated that field is exported.
//  type Application struct {
//      Config  *configuration.ApplicationConfig `autowire:""`
//  }
// Dependency injection of exported field is supported with the difference that we don't
// provide additional Setter function.
//
// Interface Injection
//
// Often it's better to rely on interface instead of concrete type, in order to
// accomplish this decoupling we should specify interfaces as a type of struct
// fields. The following snippet demonstrate that:
//  type UserService struct {
//      userRoleRepository UserRoleRepository `autowire:"repository/InMemoryUserRoleRepository"`
//  }
//  func (u *UserService) SetUserRoleRepository(userRoleRepository UserRoleRepository) {
//      u.userRoleRepository = userRoleRepository
//  }
// UserRoleRepository is simply an interface and InMemoryUserRoleRepository is a
// struct, which implements that interface. As dependency injection is executed
// on unexported field, providing Setter is required. Just to highlight unexported
// fields needs Setter while exported not. For more information take a look at
// example https://github.com/go-autowire/autowire/tree/main/example package.
// Very Simplified Example:
//		type App struct {}
//		func init()  {
//			autowire.Autowire(&App{})
//		}
// As mentioned above Autowire function should be invoked in the package init function.
func Autowire(v interface{}) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
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
	case reflect.Invalid:
		log.Panicln("invalid reflection type")
	default: // reflect.Array, reflect.Struct, reflect.Interface, etc.
		log.Panicln("unsupported: " + value.Type().String() + " value")
	}
}

// Autowired function returns fully initialized with all dependencies instance, which is ready to be used.
// As the result is empty interface, type assertions is required before using the instance.
// Take a look at https://golang.org/ref/spec#Type_assertions for more information.
// The following snippet demonstrate how could be done :
// 	 app := autowire.Autowired(app.Application{}).(*app.Application)
func Autowired(v interface{}) interface{} {
	value := reflect.ValueOf(v)
	var path string
	switch value.Kind() {
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

// Close function invoke Close method on each autowired struct
// which implements io.Closer interface, so currently active
// occupied resources (connections, channels, descriptor, etc.)
// could be released. Returning slice of occurred errors.
func Close() []error {
	log.Println("Closing...")
	var errors []error
	for _, dependency := range dependencies {
		valueDepend := reflect.ValueOf(dependency)
		closerType := reflect.TypeOf((*io.Closer)(nil)).Elem()
		if valueDepend.Type().Implements(closerType) {
			err := dependency.(io.Closer).Close()
			if err != nil {
				log.Println(err.Error())
				errors = append(errors, err)
			}
		}
	}
	return errors
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
					setValue(value, elem, i, dependency)
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
