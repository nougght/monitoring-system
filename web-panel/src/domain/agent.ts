

export interface Agent {
    id: string
    name: string
    description?: string
    createdAt: string
    status?: string
    lastSeenAt?: Date
    isOnline: boolean
}