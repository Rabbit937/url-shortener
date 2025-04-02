fetch('http://localhost:8080/api/create', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        long_url: 'https://example.com'
    })
}).then(response => {
    console.log(response)
})