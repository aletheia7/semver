#### semver 
Go package that compares semver.org version strings

#### Install 
```bash
go get github.com/aletheia7/semver
```

#### Documentation
[godoc semver](http://godoc.org/github.com/aletheia7/semver) 

[semver.org](http://semver.org/) version strings only allow [0-9A-Za-z-]. This
package allows unicode letters in place of A-Z and a-z; i.e. 
3.24.3-β+20150115102400 is acceptable. The allowance of unicode characters
makes this package noncompliant with semver.org.

![LGPL](http://www.gnu.org/graphics/lgplv3-147x51.png)
