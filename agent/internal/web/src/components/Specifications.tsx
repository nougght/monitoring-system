import { useEffect, useState } from "react";
import { getSpecs } from "../api/client/monitoringAgentAPI";
import type { Specs } from "../domain/specs";
import { convertBytesToGB } from "../util/units";

interface SpecificationsProps {

}
const Specifications = ({ }: SpecificationsProps) => {
    const [specs, setSpecs] = useState<Specs | null>(null);

    useEffect(() => {
        const fetchSpecs = async () => {
            const fetchedSpecs = await getSpecs({

            });
            if (fetchedSpecs.status == 200) {
                setSpecs(fetchedSpecs.data);
                console.log(specs);
            } else {
                console.log(fetchedSpecs.status);
            }
        }
        fetchSpecs();
    }, []);

    return (
        <div>
            <h1>Характеристики</h1>
            <section>
                <h2>Компьютер</h2>
                <p>Имя компьютера: {specs?.host?.hostName}</p>
                <p>Тип операционной системы: {specs?.host?.osType}</p>
                <p>Операционная система: {specs?.host?.os}</p>
                <p>Версия операционной системы: {specs?.host?.osVersion}</p>
                <p>Версия ядра операционной системы: {specs?.host?.osKernelVersion}</p>
                <p>Архитектура операционной системы: {specs?.host?.osArch}</p>
            </section>
            <section>
                <h2>CPU</h2>
                <p>Модель: {specs?.cpu?.modelName}</p>
                <p>Ядер: {specs?.cpu?.coreCount}</p>
                <p>Логических ядер: {specs?.cpu?.logicalCoreCount}</p>
            </section>
            <section>
                <h2>Диски</h2>
                {specs?.disk?.map((disk, i) => (
                    <div key={disk.device}>
                        <h3>Диск {disk.device}</h3>
                        <p>Файловая система: {disk.fsType}</p>
                        <p>Объем: {convertBytesToGB(disk.total).toFixed(2)} GB</p>
                    </div>
                ))}
            </section>
        </div>
    );
};

export default Specifications;