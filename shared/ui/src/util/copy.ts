
export const CopyToClipboard = (text: string) => {
    const textArea = document.createElement("textarea");
    textArea.value = text;

    textArea.style.top = "0";
    textArea.style.left = "0";
    textArea.style.position = "fixed";

    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();

    try {
        const successful = document.execCommand('copy');
        if (successful) {
            console.log('text copied!');
        } else {
            console.clear();
        }
    } catch (err) {
        console.error('failed to copy', err);
    }

    document.body.removeChild(textArea);
}