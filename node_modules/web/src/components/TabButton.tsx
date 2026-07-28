interface TabButtonProps {
    text: string;
    onClick: () => void;
    isActive: boolean;
}
const TabButton = ({text, onClick, isActive}: TabButtonProps) => {
    return (
        <button className={`tab-button ${isActive ? "active" : ""}`} onClick={onClick}>
            {text}
        </button>
    )
}

export default TabButton;