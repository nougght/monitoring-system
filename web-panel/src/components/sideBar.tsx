import { NavLink } from "react-router-dom"


export interface NavItem {
    id: string
    title: string
    iconSrc: string
    path: string
    countLabel: number
}

export interface SideBarData {
    iconSrc: string
    title: string
    items: NavItem[]
    // TODO: add more sections
}

// interface SideBarProps {
//     data: SideBarData
//     onSwitch: (id: string) => void
//     activeId: string
// }


const SideBarButton = ({ data }: { data: NavItem }) => {
    return (
        <NavLink
            key={data.path}
            to={data.path}
            className={({ isActive }) => `side-bar-button ${isActive ? "active" : ""}`}
        >
            {data.iconSrc != "" && <img src={data.iconSrc} width="15px" height="15px" />}
            <span className="side-bar-button-title">{data.title}</span>
            <span className="side-bar-button-count">{data.countLabel}</span>
        </NavLink>
    )
}

export const SideBar = ({ data}: {data: SideBarData}) => {
    return (
        <aside className="side-bar">
            <div className="side-bar-header">

            </div>
            <nav>
                {
                    Array.from(data.items.values()).map((item) => (
                        <SideBarButton
                            key={item.id}
                            data={item}
                        />
                    ))
                }
            </nav>
        </aside>
    )
}