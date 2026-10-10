import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {MutationCache, QueryClient, QueryClientProvider} from '@tanstack/react-query'
import {ReactQueryDevtools} from '@tanstack/react-query-devtools'
import {ApiError} from "./api/apiFetch.ts";

const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            // 4xx не повторять, только 5xx и сетевые ошибки
            retry: (n, e) => n < 2 && (!(e instanceof ApiError) || e.status === 0 || e.status >= 500),
        },

    },
    mutationCache: new MutationCache({
        onError: (e) => console.log(e.message),
    }),
})

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <App/>
            <ReactQueryDevtools initialIsOpen={false}/>
        </QueryClientProvider>
    </StrictMode>,
)