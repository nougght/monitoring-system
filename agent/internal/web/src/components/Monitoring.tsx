const Monitoring = ({focusedWindow}: {focusedWindow: string}) => {
    return (
        <section>
            <h2>Active window</h2>
            <p>{focusedWindow}</p>
        </section>
    )
}

export default Monitoring;