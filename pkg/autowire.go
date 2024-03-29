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
	//nolint:gochecknoglobals
	dependencies map[string]interface{}
	//nolint:gochecknoglobals
	requiredDependencies map[string]map[string]interface{}
	//nolint:gochecknoglobals
	currentProfile = internal.GetProfile()
)

// Tag respresents autowire go tag.
const Tag = "autowire"

//nolint:gochecknoinits
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Init Autowire Context")
	dependencies = make(map[string]interface{})
	requiredDependencies = make(map[string]map[string]interface{})
}

// RunProd executes function in case environment is production only, this way
// it is preventing execution of it inside go tests.
// This flexibility could help if you want to skip autowiring struct in our tests.
func RunProd(runFunc func()) {
	if currentProfile == internal.Production {
		runFunc()
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
// example package: https://github.com/go-autowire/autowire/tree/main/example.
// Very Simplified Example:
//		type App struct {}
//		func init()  {
//			Autowire(&App{})
//		}
// As mentioned above Autowire function cloud be invoked in the package init function,
// but also it is possible to do it in the main function of the application,
// or separate files, which will be responsible for autowiring all the structs.
func Autowire(values ...interface{}) {
	for _, v := range values {
		autowire(v)
		depPath := getStructPtrFullPath(reflect.ValueOf(v))
		if uncompletedDepMap, found := requiredDependencies[depPath]; found {
			for uncompleted := range uncompletedDepMap {
				if dep, ok := dependencies[uncompleted]; ok { // check for tags
					delete(uncompletedDepMap, uncompleted)
					autowireDependencies(reflect.ValueOf(dep))
				}
			}
		}
	}
}

func autowire(v interface{}) {
	value := reflect.ValueOf(v)
	switch value.Kind() { //nolint:exhaustive
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
		log.Panicf("autowiring structs is unsupported, expected to receive struct pointer(*%s)",
			value.Type().String())
	}
}

// Autowired function returns fully initialized instance with all dependencies, which is ready to be used.
// As the result is empty interface, type assertions is required before using the instance.
// Take a look at https://golang.org/ref/spec#Type_assertions for more information.
// The following snippet demonstrate simple usage of it :
// 	 app := Autowired(app.Application{}).(*app.Application)
func Autowired(v interface{}) interface{} {
	value := reflect.ValueOf(v)
	var path string
	switch value.Kind() { //nolint:exhaustive
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
	requiredDependencies = make(map[string]map[string]interface{})
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
	structType := getStructPtrFullPath(value)
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
						log.Println(msg)
					} else {
						log.Println(msg + ", ready for spy")
					}
					markStructUninitialized(structType, tag)
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
				} else {
					markStructUninitialized(structType, getStructPtrFullPath(t))
				}
			}
		}
	}
}

func markStructUninitialized(structType string, depName string) {
	if depMap, ok := requiredDependencies[depName]; ok {
		depMap[structType] = true
	} else {
		requiredDependencies[depName] = map[string]interface{}{}
		depMap = requiredDependencies[depName]
		depMap[structType] = true
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
