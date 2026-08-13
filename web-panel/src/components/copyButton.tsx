import { CopyToClipboard} from "@monitoring-system/shared"
import copyIcon from "../assets/copy.svg"

export const CopyButton = (props: { text: string }) => {
    const handleClick = () => {
        CopyToClipboard(props.text)
    }
    return (
        <button className="copyButton" onClick={handleClick}>
            <img src={copyIcon} width="16" height="16"/>
        </button>
    )
}