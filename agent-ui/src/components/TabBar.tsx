import TabButton from "./TabButton"

interface TabBarProps {
    tabs: string[];
    onSwitch: (index: number) => void;
    activeTab: number;
}
const TabBar = ({ tabs, onSwitch, activeTab }: TabBarProps) => {
    return (
        <div className="tabBar">
            {tabs.map((tab, i) => (
                <TabButton key={tab} text={tab} onClick={() => onSwitch(i)} isActive={i === activeTab} />
            ))}
        </div>
    )
}

export default TabBar;