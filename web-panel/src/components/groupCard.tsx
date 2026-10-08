import type {Group} from "../domain/group.ts";
import cardSt from "./card.module.css";
import styles from "./groupCard.module.css";


export const GroupCard = ({ group, onClick }: { group: Group, onClick: (id: string) => void }) => {
    return (
        <div className={`${cardSt.card} ${styles.groupCard}`} onClick={() => onClick(group.id)}>
            <h2>{group.name}</h2>
            <p>{group.description ?? "no description"}</p>
            <div>
                <span>0 ПК </span>
                <span>0 онлайн</span>
            </div>

        </div>
    )
}