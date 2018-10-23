import * as React from "react";

// Prevents content from scrolling
export default class NoScroll extends React.Component {

    constructor(props) {
        super(props);
        this.dom = React.createRef();
    }

    componentDidMount() {
        // Attach raw so we can set passive false to prevent scrolling
        this.dom.current.addEventListener("touchmove", NoScroll.preventDefault, {passive: false});
    }

    static preventDefault(e) {
        e.preventDefault();
    }

    componentWillUnmount() {
        this.dom.current.removeEventListener("touchmove", NoScroll.preventDefault);
    }

    render() {
        return <div ref={this.dom}>
            {this.props.children}
        </div>
    }
}
