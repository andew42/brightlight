import * as React from "react";
import './HorizontalSlider.css';

// props - min, max, pos, onPosChange(pos), sliderColour, className
export default class HorizontalSlider extends React.Component {

    constructor(props) {
        super(props);
        this.domElement = React.createRef();
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

        // Control range (typically 100)
        let range = this.props.max - this.props.min;

        let r = el.getBoundingClientRect();

        // Ensure mouse / finger is within control
        if (mouseX < r.left || mouseX > r.right || mouseY > r.bottom || mouseY < r.top) {
            return;
        }

        // Where is the mouse positioned within track
        let pos;
        if (mouseX < (r.left + padding)) {
            // Position is in the left dead zone
            pos = this.props.min;
        } else if (mouseX > (r.right - padding)) {
            // Position is in the right dead zone
            pos = this.props.max;
        } else {
            // Calculate position in pixels
            pos = mouseX - (r.left + padding);
            // Convert to value between min and max
            pos = (pos / (r.width - 2 * padding)) * range + this.props.min;
        }
        // Inform the caller
        this.props.onPosChange(pos);
    };

    doTouchUpdate(e) {
        this.updatePosition(this.domElement.current, e.targetTouches[0].clientX, e.targetTouches[0].clientY);
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
                    onTouchStart={e => this.doTouchUpdate(e)}
                    onTouchMove={e => this.doTouchUpdate(e)}
                    onMouseDown={e => this.doMouseDown(e)}
                    onMouseMove={e => this.doMouseMove(e)}
                    onMouseUp={() => this.mousedown = false}>
            <div className='hs-slider-track'
                 style={{backgroundColor: this.props.sliderColour}}>
                <div className='hs-slider-thumb-container'
                     style={{left: ((range - (range - this.props.pos)) / (range / 100)) + '%'}}>
                    <div className='hs-slider-thumb'/>
                    <div className='hs-slider-thumb-anno'>{this.props.pos.toFixed(0)}</div>
                </div>
            </div>
            <div className='hs-label'>{this.props.label}</div>
        </div>
    }
}
