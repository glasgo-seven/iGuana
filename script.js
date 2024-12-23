const examples = [
	{
		"title": "Basic Template Usage",
		"text": "Simplify your HTML by using templates for common elements like headers and footers. This reduces repetition and makes maintenance easier across your entire site.",
		"code": '{{ template "header" }}\n<div class="container">\n  <h1>Welcome</h1>\n  <p>This is a templated page</p>\n</div>\n{{ template "footer" }}'
	},
	{
		"title": "Layout Template System",
		"text": "Create a base layout that all your pages can inherit from. This ensures consistency across your site while keeping your code DRY and maintainable.",
		"code": '<!-- layout.html -->\n<!DOCTYPE html>\n<html>\n<head>\n    {{ template "meta" . }}\n    {{ template "styles" . }}\n</head>\n<body>\n    {{ template "header" . }}\n    {{ template "content" . }}\n    {{ template "footer" . }}\n    {{ template "scripts" . }}\n</body>\n</html>'
	},
	{
		"title": "Dynamic Components",
		"text": "Build reusable components with conditional rendering and dynamic data. Perfect for cards, buttons, and other repeatable elements throughout your site.",
		"code": '<!-- components/card.html -->\n{{ define "card" }}\n<div class="card">\n    <img src="{{ .image }}" alt="{{ .title }}">\n    <h3>{{ .title }}</h3>\n    <p>{{ .description }}</p>\n    {{ if .button }}\n        <a href="{{ .button.url }}" class="btn">{{ .button.text }}</a>\n    {{ end }}\n</div>\n{{ end }}'
	}
]

window.onload =  function() {
	var target = document.getElementById("content");

	examples.forEach(ex => {
		target += `<h4>${ex.title}</h4><p>${ex.text}</p><span>${ex.code}</span>`;
	});
}
