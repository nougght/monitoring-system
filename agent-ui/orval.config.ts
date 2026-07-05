export default {
    agent: {
        input: 'http://127.0.0.1:8088/swagger/doc.json',
        output: {
            mode: 'split',
            schemas: 'src/api/models',
            target: 'src/api/client',
            client: 'fetch',
            baseUrl: 'http://127.0.0.1:8088',
        },
    },
};