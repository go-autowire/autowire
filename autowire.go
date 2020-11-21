package autowire

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"regexp"
	"strings"
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

func RunOmitTest(f func()) {
	if currentProfile != _Testing {
		f()
	}
}

func Autowire(v interface{}) {
	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Invalid:
		log.Println("invalid reflection type")
	case reflect.Ptr:
		structType := getStructPtrFullPath(value)
		_, ok := dependencies[structType]
		if ok {
			log.Printf("%s already autowired. Going go overwrite it`s value!", structType)
		} else {
			autowireDependencies(value)
			dependencies[structType] = v
		}
	default: // reflect.Array, reflect.Struct, reflect.Interface
		log.Println(value.Type().String() + " value")
	}
	log.Println(dependencies)
}

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
	log.Println(path)
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
			log.Println("......." + field.Type.String())
			if tag != "" {
				log.Printf("field type: %v\n", getFullPath(field.Type.PkgPath(), field.Type.String()))
				currentDep := findDependency(tag)
				if len(currentDep) == 0 {
					msg := "Unknown dependency " + tag + " found none"
					if currentProfile != _Testing {
						log.Panicln(msg)
					} else {
						log.Println(msg)
					}
				} else {
					v := reflect.ValueOf(currentDep[0])
					if v.Type().Implements(field.Type) {
						t = reflect.New(v.Type())
						log.Println("found dependency by tag " + tag + " will be used " + t.Type().String())
						elem.Field(i).Set(reflect.ValueOf(currentDep[0]))
					} else {
						log.Panicln(v.Type().String() + " doesnt Implements: " + field.Type.String())
					}
				}
			} else {
				t = reflect.New(elem.Type().Field(i).Type.Elem())
				log.Printf("field type: %v", getStructPtrFullPath(t))
				dependency, found := dependencies[getStructPtrFullPath(t)]
				if found {
					elem.Field(i).Set(reflect.ValueOf(dependency))
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

type AutoClosable interface {
	Close()
}

func Run(mainFunction func()) {
	defer close()
	mainFunction()
}

func close() {
	println("HEREE")

}
