class Hello extends React.Component {
  render() {
    let who = "World";
    return (
      <h1> Hello {who}! </h1>
    );
  }
}

ReactDOM.render( <Hello/>, document.querySelector("#root"));
