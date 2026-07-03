import { useEffect, useState } from "react";
import { getSpecs } from "../api/specs";
import type { CpuSpecs } from "../domain/specs";

interface SpecificationsProps {

}
const Specifications = ({}: SpecificationsProps) => {
    const [specs, setSpecs] = useState<CpuSpecs | null>(null);

    useEffect(() => {
        const fetchSpecs = async () => {
            const fetchedSpecs = await getSpecs();
            try {
                setSpecs(fetchedSpecs);
                console.log(specs);
            } catch (error) {
                console.error("Error fetching specs:", error);
            }
        }
        fetchSpecs();
    }, []);
    
    return (
        <div>
        <h1>Характеристики</h1>
        <section>
            <h2>CPU</h2>
            <p>Модель: {specs?.modelName}</p>
            <p>Ядер: {specs?.coreCount}</p>
            <p>Логических ядер: {specs?.logicalCoreCount}</p>
        </section>
        </div>
    )
}

export default Specifications;