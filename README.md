# iGuana

iHTML (interpreted Hyper Text Markdown Language) is a super-set of HTML


#### TODO
- [ ]	File validation
	- [ ]	html tags
	- [ ]	http syntax
- [ ]	Versions
- [ ]	[OOSS](./OOSS.md)

It introduces some new concepts:

- Hyper Text Template (HTT)  
	HTT is a self-sustaining block of HTML code, that can be safely inserted into any web-page and work correctly

- Hyper Text Template Package (HTTP)  
	HTTP is a collection of HTTs, used to import them into HTML file


There are also new tags:

- `<import src='package.http' />`  
	import all the templates in the package file

- `<require template='template.htt'>`  
	imports only one specific template by file name

- `<call template='templateName' />`  
	insert the template contents on its place

---

**template_A.htt**
```html
<span>
	Lorem ipsum dolor sit amen
</span>
```

**template_B.htt**
```html
<require template='template_A.htt' />
<div>
	<h1>Hello World!</h1>
	<call template='templateA' />
</div>
```

**template_C.htt**
```html
<h2>
	I am being ignored :(
</h2>
```

**package.http**
```cs
[PackageName]
// Silenced, can be accessed only by other templates
~ template_A.htt [PackageName]
  
// Usual import
+ template_B.htt [HelloWorld]

// Permanent exclusion
! template_C.htt []
```

**index.html**
```html
<import src='package.http' />
<html>
	<body>
		<call template='PackageName.HelloWorld' />
	</body>
</html>
```

---
