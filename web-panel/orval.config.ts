export default {
    web: {
        input: 'http://127.0.0.1:8091/api/v1/swagger/doc.json',
        output: {
            mode: 'split',
            schemas: 'src/api/models',
            target: 'src/api/client',
            client: 'fetch',
            baseUrl: 'http://127.0.0.1:8091/api/v1',
        },
    },
};