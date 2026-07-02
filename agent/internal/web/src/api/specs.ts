import type { CpuSpecs } from "../domain/specs"

export const getSpecs = async (): Promise<CpuSpecs> => {
    try {
        const URL = "http://127.0.0.1:8088/specs"
        const response = await fetch(URL)
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`)
        }
        const data = await response.json() as CpuSpecs
        console.log(data);
        return data
    } catch (error) {
        console.error("Error fetching specs:", error);
        throw error;
    }
}