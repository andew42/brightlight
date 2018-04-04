import * as React from "react";
import * as ReactDOM from "react-dom";

export class Bouncefix extends React.Component {
    constructor(props) {
        super(props);
        this.onTouchEnd = this.onTouchEnd.bind(this);
        this.onTouchStart = this.onTouchStart.bind(this);
        this.onTouchMove = this.onTouchMove.bind(this);
    }

    scrollToEnd(el) {
        let curPos = el.scrollTop,
            height = el.offsetHeight,
            scroll = el.scrollHeight;

        // If at top, bump down 1px
        if (curPos <= 0) {
            el.scrollTop = 1;
        }

        // If at bottom, bump up 1px
        if (curPos + height >= scroll) {
            el.scrollTop = scroll - height - 1;
        }
    }

    onTouchStart() {
        let el = ReactDOM.findDOMNode(this);
        let isScrollable = el.scrollHeight > el.offsetHeight;

        // If scrollable, adjust
        if (isScrollable) {
            this._blockTouchMove = false;
            return this.scrollToEnd(el);
        }
        // Else block touchmove
        else {
            this._blockTouchMove = true;
        }
    }

    onTouchMove(e) {
        if (this._blockTouchMove) {
            e.preventDefault();
        }
    }

    onTouchEnd() {
        this._blockTouchMove = false;
    }

    render() {
        return <div onTouchStart={this.onTouchStart}
                    onTouchMove={this.onTouchMove}
                    onTouchEnd={this.onTouchEnd}
                    onTouchCancel={this.onTouchEnd}>
            {this.props.children}
        </div>
    }
}
