import * as React from "react";
import './VerticalSlider.css';

// props - min, max, pos, onPosChange(pos), sliderColour, className
export default class VerticalSlider extends React.Component {

    constructor(props) {
        super(props);
        this.domElement = React.createRef();
        this.doTouchUpdate = this.doTouchUpdate.bind(this);
    }

    componentDidMount() {
        // Attach raw so we can set passive false to prevent scrolling
        this.domElement.current.addEventListener("touchstart", this.doTouchUpdate, {passive: false});
        this.domElement.current.addEventListener("touchmove", this.doTouchUpdate, {passive: false});
    }

    componentWillUnmount() {
        this.domElement.current.removeEventListener("touchstart", this.doTouchUpdate);
        this.domElement.current.removeEventListener("touchmove", this.doTouchUpdate);
    }

    // Helper converts pixel positions into range position and calls onValueChange
    updatePosition(el, mouseX, mouseY) {

        // We expect the element to have some padding which will be
        // treated as a dead zone to allow easy selection of min max
        // values. Padding is expected to be in pixels
        let padding = parseInt(window.getComputedStyle(el, null).padding, 10);
        if (padding === undefined) {
            padding = 0;
        }

        // Control range (typically 360 or 100)
        let range = this.props.max - this.props.min;

        let r = el.getBoundingClientRect();

        // Ensure mouse / finger is within control
        if (mouseX < r.left || mouseX > r.right || mouseY > r.bottom || mouseY < r.top) {
            return;
        }

        // Where is the mouse positioned within track
        let pos;
        if ((r.bottom - mouseY) < padding) {
            // Position is in the lower dead zone
            pos = this.props.min;
        } else if (mouseY < (r.bottom - r.height + padding)) {
            // Position is in the upper dead zone
            pos = this.props.max;
        } else {
            // Calculate position in pixels
            pos = (r.bottom - padding) - mouseY;
            // Convert to value between min and max
            pos = (pos / (r.height - 2 * padding)) * range + this.props.min;
        }
        // Inform the caller
        this.props.onPosChange(pos);
    };

    doTouchUpdate(e) {
        if (this.domElement === undefined)
            return;
        this.updatePosition(this.domElement.current, e.targetTouches[0].clientX, e.targetTouches[0].clientY);
        e.preventDefault();
    }

    doMouseDown(e) {
        if (!('ontouchstart' in window)) {
            this.mousedown = true;
            this.updatePosition(this.domElement.current, e.clientX, e.clientY);
        }
    }

    doMouseMove(e) {
        if (this.mousedown) {
            this.updatePosition(this.domElement.current, e.clientX, e.clientY);
        }
    }

    render() {
        let range = this.props.max - this.props.min;
        return <div ref={this.domElement}
                    className={this.props.className}
                    onMouseDown={e => this.doMouseDown(e)}
                    onMouseMove={e => this.doMouseMove(e)}
                    onMouseUp={() => this.mousedown = false}>
            <div className='vs-slider-track'
                 style={{backgroundColor: this.props.sliderColour}}>
                <div className='vs-slider-thumb-container'
                     style={{top: ((range - this.props.pos) / (range / 100)) + '%'}}>
                    <div className='vs-slider-thumb'/>
                    <div className='vs-slider-thumb-anno'>{this.props.pos.toFixed(0)}</div>
                </div>
            </div>
        </div>
    }
}
