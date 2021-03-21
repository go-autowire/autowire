// Package autowire is a framework that makes easy using the dependency
// injection in Golang.
//
// The main purpose of using Autowire as dependency injection framework is
// to eliminate usage of globals without the tedious approach to manually
// wiring all the dependencies together. Current approach of Autowire
// framework is relying on dependency injection via using struct tags and
// reflection. All the dependencies are injected via golang reflection.
//
// Basic usage is explained in the package-level example below. If you're new
// to Autowire, start there!
//
// Testing of Autowire Applications
//
// To write unit or end-to-end tests of your application, you can use functions
// provided by atesting package. For more information take a look at
// https://godoc.org/github.com/go-autowire/autowire/atesting.
package autowire
