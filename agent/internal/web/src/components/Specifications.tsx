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
                <h3>{specs?.cpu?.modelName}</h3>
                <p>Производитель: {specs?.cpu?.manufacturer}</p>
                <p>Архитектура: {specs?.cpu?.architecture}</p>
                <p>Сокет: {specs?.cpu?.socketDesignation}</p>
                <p>Количество ядер: {specs?.cpu?.numberOfCores}</p>
                <p>Количество доступных ядер: {specs?.cpu?.numberOfEnabledCore}</p>
                <p>Количество логических ядер: {specs?.cpu?.numberOfLogicalProcessors}</p>
                <p>Размер L2 кэша: {specs?.cpu?.l2CacheSize == 0 ? 'Неизвестно' : `${specs?.cpu?.l2CacheSize / 1024} Mb`}</p>
                <p>Размер L3 кэша: {specs?.cpu?.l3CacheSize == 0 ? 'Неизвестно' : `${specs?.cpu?.l3CacheSize / 1024} Mb`}</p>
                <p>Состояние: {specs?.cpu?.availability}</p>
                <p>Текущая частота процессора: {specs?.cpu?.currentClockSpeed}</p>
                <p>Максимальная частота процессора: {specs?.cpu?.maxClockSpeed}</p>
                <p>Идентификатор процессора: {specs?.cpu?.processorId}</p>
                <p>Степпинг: {specs?.cpu?.stepping}</p>
                <p>Разрядность процессора: {specs?.cpu?.dataWidth}</p>
                <p>Виртуализация: {specs?.cpu?.virtualizationFirmwareEnabled ? 'Включена' : 'Выключена'}</p>
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
            <section>
                <h2>Оперативная память</h2>
                <p>Общий объем: {specs?.memory?.total}</p>
                {specs?.memory?.physicalMemoryList?.map((mem, i) => (
                    <div key={i}>
                        <h3>{mem.deviceLocator}</h3>
                        <p>Тип памяти: {mem.memoryType?.toString()}</p>
                        <p>Форм-фактор: {mem.formFactor?.toString()}</p>
                        <p>Объем: {convertBytesToGB(mem.capacity).toFixed(2)} GB</p>
                        <p>Поддерживаемя частота: {mem.speed}</p>
                        <p>Настроенная частота: {mem.configuredClockSpeed}</p>
                        <p>Производитель: {mem.manufacturer}</p>
                        <p>Модель: {mem.modelName}</p>
                        <p>Серийный номер: {mem.serialNumber}</p>
                        <p>Подключение: {mem.bankLabel}</p>
                        <p>Можно менять память без выключения системы: {mem.hotSwappable ? 'Да' : 'Нет'}</p>
                        <p>Можно ли вынимать память: {mem.removable ? 'Да' : 'Нет'}</p>
                        <p>Можно ли заменять память: {mem.replaceable ? 'Да' : 'Нет'}</p>
                    </div>
                ))}
            </section>
        </div>
    );
};

export default Specifications;