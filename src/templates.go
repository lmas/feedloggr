package feedloggr2

const HTML_BODY = `
<!doctype html>
<html>
	<head>
		<meta charset="utf-8">
		<title>
			{{format .Date}} | News
		</title>
		<style type="text/css">
			body {
				margin: 0px;
			}
			#content {
				background-color: #FFF;
				color: #444;
				font-family: monospace;
				font-size: 14px;
				padding: 25px;
				padding-bottom: 0px;
			}
			ol {
				margin: 0px;
			}
			li.item {
				margin-bottom:3px;
			}
			a {
				text-decoration:none;
				color: #444;
			}
			a:hover {
				text-decoration:none;
				color: #000;
			}
			a:visited {
				color: #AAA;
			}
			.center {
				text-align: center;
			}

			#date_nav {
				text-align: center;
				margin-bottom: 25px;
			}
			#date_nav>a {
				border: 1px solid #DDD;
				border-radius: 3px;
				padding: 6px 12px;
				line-height: 1.4;
				font-weight: bold;
			}

			.panel {
				border: 1px solid #DDD;
				background-color: #FFF;
				border-radius: 3px;
				margin-bottom: 20px;
			}
			.panel-heading {
				background-color: #EEE;
				border-top-left-radius: 3px;
				border-top-right-radius: 3px;
				border-bottom: 1px solid #DDD;
				padding: 10px 15px;
			}
			.panel-heading>a:visited {
				color: #444;
			}
			.panel-heading>a:hover {
				color: #000;
			}
			.panel-body {
				padding: 15px;
				line-height: 1.4;
			}
		</style>
	</head>

	<body>
		<div id="content">
			<div id="date_nav">
				<a class="btn" href="{{date_link 24 .Date}}">&lt;</a>
				<a class="btn" href="index.html">Latest</a>
				<a class="btn" href="{{date_link -24 .Date}}">&gt;</a>
			</div>
			{{ range .Feeds }}
				<div class="panel panel-default">
					<div class="panel-heading">
						<a href="{{.Url}}" rel="nofollow">{{.Title}}</a>
					</div>
					<div class="panel-body">
						<ol>{{ range .Items }}
							<li class="item"><a href="{{.Url}}" rel="nofollow" class="item">
								{{.Title}}
							</a></li>
						{{ end }}</ol>
					</div>
				</div>
			{{ end }}
		<p class="center">
			Generated with Feedloggr2
		</p>
		</div>
	</body>
</html>
`
