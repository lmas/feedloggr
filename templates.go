package feedloggr

const tmplCSS string = `
body {
	margin: 30px 5%;
	line-height: 1.5;
	font-size: 13px;
	font-family: monospace;
	background-color: #FFF;
	color: #444;
}
a, a:visited {
	color: #444;
	text-decoration: none;
}
a:hover {
	color: #000;
	text-decoration: underline;
}

nav {
	text-align: center;
	margin-bottom: 20px;
}
nav > a:hover {
	background-color: #E6E6E6;
	border-color: #ADADAD;
}
section {
	background-color: #FFF;
	margin-bottom: 20px;
}
section > h1, nav > a {
	border: 1px solid #DDD;
	border-radius: 3px;
	padding: 10px;
	margin: 0;
	font-size: 16px;
	font-weight: bold;
	text-align: center;
}
section > h1 {
	background-color: #EEE;
}
section > ul, section > p {
	padding: 0;
	margin: 20px 0;
	list-style: none;
}
section > ul > li {
        margin-bottom: 5px;
}
section > ul > a:visited {
	color: #AAA;
}
section > ul > hr {
        padding: 0;
        border: 0;
        height: 1px;
        margin: 5px 0px;
        background-image: linear-gradient(to right, #ccc, #fff);
}
footer {
	text-align: center;
	font-size: 12px;
}
`

const tmplPage string = `
<!doctype html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>{{.CurrentDate}} | Feedloggr</title>
	<link href="./style.css" rel="stylesheet" type="text/css">
</head>
<body>
	<nav>
		<a href="{{.NextDate}}.html">&lt;</a>
		<a href="index.html">Latest</a>
		<a href="{{.PrevDate}}.html">&gt;</a>
	</nav>
	{{range .Feeds}}
	<section>
		<h1><a href="{{.URL}}" rel="nofollow">{{.Title}}</a></h1>
		{{if .Error}}
		<p>Error while updating feed:<br />{{.Error}}</p>
		{{else}}
		<ul>
		{{- range $i, $row := .Items}}
			{{- if $i}}<hr>{{end}}
			<a href="{{$row.URL}}" rel="nofollow"><li>{{$row.Title}}</li></a>
		{{- end}}
		</ul>
		{{end}}
	</section>
	{{else}}
	<p class="center">Sorry, no news for today!</p>
	{{end}}
	<footer>
		Generated with <a href="https://github.com/lmas/feedloggr">Feedloggr</a>
	</footer>
</body>
</html>
`
