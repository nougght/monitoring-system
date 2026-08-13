

export interface Agent {
    id: string
    name: string
    description?: string
    createdAt: string
    lastSeenAt?: Date
    isOnline: boolean
}