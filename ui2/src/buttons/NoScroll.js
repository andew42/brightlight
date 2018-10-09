import * as React from "react";

// Prevents content from scrolling
export default class NoScroll extends React.Component {

    constructor(props) {
        super(props);
        this.dom = React.createRef();
    }

    componentDidMount() {
        // Attach raw so we can set passive false to prevent scrolling
        if (this.dom.current !== null) {
            this.dom.current.addEventListener("touchmove", this.preventDefault, {passive: false});
        }
    }

    preventDefault(e) {
        e.preventDefault();
    }

    componentWillUnmount() {
        if (this.dom.current !== null)
            this.dom.current.removeEventListener("touchmove", this.preventDefault);
    }

    render() {
        return <div ref={this.dom}>
            {this.props.children}
        </div>
    }
}
