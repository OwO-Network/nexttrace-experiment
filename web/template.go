package web

import (
	"io/ioutil"

	"github.com/OwO-Network/nexttrace-enhanced/config"
)

func writeTemplateFile() error {
	var err error
	var path string
	path, err = config.ConfigFromUserHomeDir()
	if err != nil {
		path, err = config.ConfigFromRunDir()
		if err != nil {
			return err
		}
	}

	content := `{{ define "index.tmpl" }}
	<!DOCTYPE html>
	<html>
	
	<head>
		<title>{{.title}}</title>
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
			integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
		<style>
			html {
				height:100%;
			}
	
			body {
				display: flex;
				text-align: center;
			}
	
			.IPInput {
				margin-top: 20px;
			}
	
			.results {
				width: 100;
				margin: auto;
				margin-top: 10px;
			}
	
			.ipinput {
				margin-top: 30px;
				display: flex;
				justify-content: center;
				text-align: center;
			}
	
			table {
				text-indent: 0;
				border-color: inherit;
				border-collapse: collapse;
			}
	
			input {
				padding: 5px 10px;
				border-radius: 5px;
				border: 1px solid #bfbfbf;
				margin-right: 30px;
				width: 180px;
			}
	
			.btn {
				font-size: 14px;
				min-width: 100px;
				border-radius: 20px;
				vertical-align: middle;
			}
	
			tbody {
				--tw-divide-opacity: 1;
			}
		</style>
	</head>
	
	<body>
		<div class="container">
			<div class="IPInput">
				<h4>Traceroute 工具</h4>
				<div class="form-inline ipinput">
					<form>
						<span>IP地址:</span>
						<input type="text" id="ip" name="ip" placeholder="Enter IP Address">
	
						<input type="text" id="token" name="token" value="{{.token}}" hidden>
	
						<select id="method" name="method" class="form-control">
							<option value="icmp">ICMP</option>
							<option value="tcp">TCP</option>
						</select>
						<button type="button" id="trace" class="btn btn-primary" onclick="getData(); return false;">查询</button>
					</form>
				</div>
	
			</div>
			<div class="results">
				<table id="table" border="1" cellspacing="0" cellpadding="0" class="table table-bordered table-striped">
					<tr>
						<th>Hops</th>
						<th>IP</th>
						<th>rDNS</th>
						<th>Latency</th>
						<th>ASN</th>
						<th>Geography</th>
						<th>ISP</th>
					</tr>
				</table>
				<p id="inprogress" style="display: none;"></p>
	
			</div>
			<script>
				function getData() {
					cleanTable();
					document.getElementById('inprogress').innerHTML = "路由测试中...";
					document.getElementById("trace").disabled = true;
					document.getElementById("inprogress").style.display = "block";
					fetch('trace?ip=' + document.getElementById('ip').value + '&method=' + document.getElementById('method').value + '&token=' + document.getElementById('token').value)
						.then(function (res) {
							return res.json();
						})
						.then(function (myJson) {
							if (myJson == null) {
								if (document.getElementById("method").value == "tcp") {
									document.getElementById('inprogress').innerHTML = "输入的IP不正确，TCP 模式仅支持 IPv4";
								} else {
									document.getElementById('inprogress').innerHTML = "输入的IP不正确";
								}
							} else {
							
								if (myJson.Hops == null) {
									document.getElementById('inprogress').innerHTML = "数据包发送失败，宿主机可能不支持 IPv6";
								} else {
									generateTable(myJson);
									document.getElementById("inprogress").style.display = "none";
								}
							}
							
							document.getElementById("trace").disabled = false;
						});
				}
	
				function cleanTable() {
					var table = document.getElementById('table');
					while (table.rows.length > 1) {
						table.deleteRow(1);
					}
				}
	
				function generateTable(data) {
					var data = data.Hops
	
					var table = document.getElementById("table");
					for (var i = 0; i < data.length; i++) {
						if (data[i][0].Success) {
							var row = table.insertRow(table.rows.length);
							var c1 = row.insertCell(0);
							c1.innerHTML = i + 1;
							var c2 = row.insertCell(1);
							c2.innerHTML = data[i][0].Address.IP;
							var c3 = row.insertCell(2);
							c3.innerHTML = data[i][0].Hostname;
							var c4 = row.insertCell(3);
							c4.innerHTML = (data[i][0].RTT / 1000000).toFixed(2) + "ms";
							var c5 = row.insertCell(4);
							c5.innerHTML = data[i][0].Geo.Asnumber;
							var c6 = row.insertCell(5);
							c6.innerHTML = data[i][0].Geo.Country + " " + data[i][0].Geo.Prov + " " + data[i][0].Geo.City;
							var c7 = row.insertCell(6);
							c7.innerHTML = data[i][0].Geo.Owner;
						} else {
							var row = table.insertRow(table.rows.length);
							var c1 = row.insertCell(0);
							c1.innerHTML = i + 1;
							var c2 = row.insertCell(1);
							c2.innerHTML = "*"
							var c4 = row.insertCell(2);
							var c4 = row.insertCell(3);
							var c4 = row.insertCell(4);
							var c4 = row.insertCell(5);
							var c4 = row.insertCell(6);
						}
	
					}
				}
			</script>
	</body>
	
	</html>
	{{ end }}
	`

	if err = ioutil.WriteFile(path+"index.tmpl", []byte(content), 0644); err != nil {
		return err
	}

	return nil
}
