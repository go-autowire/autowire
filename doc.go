// Package Autowire is a framework that makes dependency injection in Golang.
//
// Autowire applications use dependency injection framework to eliminate usage
// of globals without the tedious approach to manually wiring all the
// dependencies together. Current approach of Autowire framework to accomplish
// dependency injection is via using struct tags.
//
// Basic usage is explained in the package-level example below. If you're new
// to Autowire, start there!
//
// Testing of Autowire Applications
//
// To write end-to-end tests of your application, you can use functions provided
//by atesting package https://godoc.org/github.com/go-autowire/autowire/atesting
package autowire
