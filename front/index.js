const socket = new WebSocket("http://127.0.0.1:12312", ["kripto"])

socket.addEventListener("open", (event) => {
	socket.send(
		'{\
			"type": "kripto",\
			"payload": [\
				{\
					"n1": 1,\
					"operation": "+",\
					"n2": 1\
				},\
				{\
					"n1": 1,\
					"operation": "+",\
					"n2": 1\
				},\
				{\
					"n1": 2,\
					"operation": "+",\
					"n2": 2\
				}\
			]\
		}'
	)
});

socket.addEventListener("message", event => {
	tit.innerHTML = event.data
	let obj = JSON.parse(event.data)
	console.log(obj);
});



