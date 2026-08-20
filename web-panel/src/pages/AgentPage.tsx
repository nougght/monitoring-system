
import TabBar from '../../../shared/ui/src/components/TabBar';
import Specifications from '../../../shared/ui/src/components/Specifications';
import { useSpecs } from "../hooks/useGetSpecs";
import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

interface Tab {
    text: string;
    content: React.ReactNode;
}

export const AgentPage = () => {
    const { id } = useParams()
    const [activeTab, setActiveTab] = useState(0)
    
    const [warning, setWarning] = useState<string | null>()
    const {
        data: specs,
        isPending: isSpecsPending,
        isError: _isSpecsError,
        error: _specsError,
        isFetching: _isSpecsFetching,
    } = useSpecs(id ?? "");


    useEffect(() => {
        if (specs?.error != null) {
            console.error(specs?.error)
            setWarning(`ошибка:${specs?.error.status} ${specs?.error.message}`)
        }
    }, [specs]);

    if (isSpecsPending) {
        return <div>Загрузка...</div>;
    }
    
    const tabs: Tab[] = [
        {
            text: "Характеристики",
            content:  specs?.specs != null ?
             <div>
                <Specifications specs={specs.specs}/>
            </div> : 
            specs?.error?.status == 404 && <div>Характеристики не найдены</div>
        }
    ]
    return (
        <div className="agent-page">
            <TabBar tabs={tabs.map((tab) => tab.text)} onSwitch={setActiveTab} activeTab={activeTab} />
            <div>
                {tabs[activeTab].content}
            </div>
            {
                warning != null &&
                <div>
                    <p>{warning}</p>
                </div>
            } 
        </div>
    )
}