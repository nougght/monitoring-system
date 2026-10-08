// import { Link, useNavigate } from "react-router-dom"
import {useState} from "react";
// import { GroupCard } from "../components/groupCard";
import styles from "./groupsPage.module.css"
// import TabBar from "../../../shared/ui/src/components/TabBar.tsx";
import {Dialog} from "radix-ui";
import {useGroups} from "../hooks/useGetGroups.ts";
import {GroupCard} from "../components/groupCard.tsx";
import {useCreateGroup} from "../hooks/useCreateGroup.ts";

interface GroupInfo {
    Name: string
    Description?: string
    Agents?: string[]
}

export const GroupsPage = () => {
    const [groupInfo, setGroupInfo] = useState<GroupInfo>({
        Name: ""
    })

    const {
        data: groups,
        isPending: isGroupsPending,
        isError: _isGroupsError,
        error: _groupsError,
        isFetching: _isGroupsFetching,
    } = useGroups();

    const {
        mutate: createGroup,
        isPending: isCreatePending,
        isError: _isCreateError,
        isSuccess: _isCreateSuccess,
        error: _createError,
    } = useCreateGroup()

    // let navigate = useNavigate()
    //
    //
    // useEffect(() => {
    //     if (groups?.error != null) {
    //         setWarning(`ошибка:${groups?.error.status} ${groups?.error.message}`)
    //     }
    // }, [groups]);
    //
    // if (isGroupsPending) {
    //     return <div>Загрузка...</div>;
    // }
    // const tabs: Tab[] = [
    //     {
    //         text: "Обзор",
    //         content:<div>
    //             <div className={styles.groupCardsContainer}>
    //                 {/*{groups?.groups != null && groups?.groups.length > 0 &&*/}
    //                 {/*    groups?.groups?.filter((a) => a.status != null)*/}
    //                 {/*        .sort((a, b) => {*/}
    //                 {/*            if (a.isOnline == b.isOnline)*/}
    //                 {/*                return 0*/}
    //                 {/*            if (a.isOnline) {*/}
    //                 {/*                return -1*/}
    //                 {/*            }*/}
    //                 {/*            return 1*/}
    //                 {/*        })*/}
    //                 {/*        .map((group) =>*/}
    //                 {/*                // {agent.status != null &&*/}
    //                 {/*                <div key={group.id}>*/}
    //                 {/*                    <GroupCard group={group} onClick={(id) => { navigate(id) }} />*/}
    //                 {/*                </div>*/}
    //                 {/*            // }*/}
    //                 {/*        )*/}
    //                 {/*}*/}
    //             </div>
    //         </div>
    //     },
    //     {
    //         text: "Управление",
    //         content: <div></div>
    //     }
    // ]

    const handleCreate = () => {
        createGroup(
            {
                name: groupInfo.Name,
                description: groupInfo.Description,
                agentIDs: groupInfo.Agents
            },
        )
    }
    return (
        <div className={styles.agentsPage}>
            <h1>Группы</h1>
            <main>
                {/*<TabBar tabs={tabs.map((tab) => tab.text)} onSwitch={setActiveTab} activeTab={activeTab} />*/}
                {isGroupsPending && <a>Загрузка</a>}
                {
                    groups && groups.map((g) =>
                        <GroupCard key={g.id} group={g} onClick={() => {
                        }}/>
                    )
                }
                <Dialog.Root>
                    <Dialog.Trigger className="btn">Новая группа</Dialog.Trigger>
                    <Dialog.Portal>
                        <Dialog.Overlay className={styles.overlay}/>
                        <Dialog.Content className={styles.content} asChild>
                            <div>
                                <div>
                                    <label className={styles.formLabel} htmlFor="name">Название</label>
                                    <input type="text" id="name" name="name" value={groupInfo?.Name}
                                           required onChange={e =>
                                        setGroupInfo({
                                            ...groupInfo,
                                            Name: e.target.value
                                        })}/>
                                    <label className={styles.formLabel} htmlFor="description">Описание</label>
                                    <input type="text" id="description" name="description"
                                           value={groupInfo?.Description}
                                           onChange={e =>
                                               setGroupInfo({
                                                   ...groupInfo,
                                                   Description: e.target.value
                                               })}/>
                                    {isCreatePending && <a>Загрузка...</a>}
                                    <button onClick={handleCreate}>Создать</button>
                                </div>
                            </div>
                        </Dialog.Content>
                    </Dialog.Portal>
                </Dialog.Root>
            </main>
            {
                _isCreateError &&
                <div>
                    <p>{_createError.message}</p>
                </div>
            }
        </div>
    )
}