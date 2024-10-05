# OOSS

### Object Oriented Style Sheet

For more convenient CSS reuse

Before:

- create a file where you describe all classes, ids and tags
- if you want to change something - your need to create a new use block and add all the styles

Now:

- Create an oop-styled file, where you can create a javascript-like class

```js
class Base {
	#padding = 8px; // private field
	margin = 8px; // public field
}

class Style extends Base {
	// padding is inherited from Base class
	margin = 16px; // new value to field
	color = #123456; // new field
}

class Error extends Style {
	padding = 1000px;   // will result in error
						// padding is a private field
}
```

Private field exist so child classes cannot change it

