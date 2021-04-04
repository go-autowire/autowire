package pkg

import (
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-autowire/autowire/pkg/internal"
)

var (
	dependencies   map[string]interface{}
	currentProfile = internal.GetProfile()
)

// Tag name
const Tag = "autowire"

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Init Autowire Context")
	dependencies = make(map[string]interface{})
}

// InitProd executes function only in production so it is preventing execution in our go tests.
// This flexibility could help if we want to skip autowiring struct in our tests.
func InitProd(initFunc func()) {
	if currentProfile == internal.Production {
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
// Following snippet shows injecting dependencies inside private structure
// fields using `autowire:""` tag:
//  type Application struct {
//      config  *configuration.ApplicationConfig `autowire:""`
//  }
// Exported field
// The use of upper-case names of the struct fields indicated that field is exported.
//  type Application struct {
//      Config  *configuration.ApplicationConfig `autowire:""`
//  }
// It is important to note that only struct pointers are supported.
//
// Interface Injection
//
// Often it's better to rely on interface instead of concrete type, in order to
// accomplish this decoupling we should specify interfaces as a type of struct
// fields. The following snippet demonstrate that:
//  type UserService struct {
//      userRoleRepository UserRoleRepository `autowire:"repository/InMemoryUserRoleRepository"`
//  }
// UserRoleRepository is simply an interface and InMemoryUserRoleRepository is a
// struct, which implements that interface. For more information take a look at
// example https://github.com/go-autowire/autowire/tree/main/example package.
// Very Simplified Example:
//		type App struct {}
//		func init()  {
//			Autowire(&App{})
//		}
// As mentioned above Autowire function should be invoked in the package init function,
// but also it is possible to do it in the main function of the application,
// or separate files, which autowire all the structs.
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
// 	 app := Autowired(app.Application{}).(*app.Application)
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
// Close functions cleans the dependency graph.
func Close() []error {
	log.Println("Closing...")
	var errors []error
	for key, dependency := range dependencies {
		valueDepend := reflect.ValueOf(dependency)
		closerType := reflect.TypeOf((*io.Closer)(nil)).Elem()
		if valueDepend.Type().Implements(closerType) {
			err := dependency.(io.Closer).Close()
			if err != nil {
				log.Println(err.Error())
				errors = append(errors, err)
			}
		}
		delete(dependencies, key)
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
					if currentProfile != internal.Testing {
						log.Panicln(msg)
					} else {
						log.Println(msg + ", ready for spy")
					}
				} else {
					v := reflect.ValueOf(currentDep[0])
					if v.Type().Implements(field.Type) {
						t = reflect.New(v.Type())
						dependency := currentDep[0]
						internal.SetFieldValue(elem, i, dependency)
					} else {
						log.Panicln(v.Type().String() + " doesnt Implements: " + field.Type.String())
					}
				}
			} else {
				t = reflect.New(elem.Type().Field(i).Type.Elem())
				dependency, found := dependencies[getStructPtrFullPath(t)]
				if found {
					internal.SetFieldValue(elem, i, dependency)
				}
			}
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
