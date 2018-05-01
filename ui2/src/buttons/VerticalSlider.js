import * as React from "react";
import './VerticalSlider.css';

// props - min, max, pos, onPosChange(pos), sliderColour, className
export default class VerticalSlider extends React.Component {

    componentDidMount() {
        this.init(this.domElement, this.props.onPosChange, this.props.min, this.props.max);
    }

    // Touch slider control (adjusts value based on position on track)
    init(el, onValueChange, min, max) {

        // We expect the element to have some padding which will be
        // treated as a dead zone to allow easy selection of min max
        // values. Padding is expected to be in pixels
        let padding = parseInt(window.getComputedStyle(el, null).padding, 10);
        if (padding === undefined) {
            padding = 0;
        }

        // Control range (typically 360 or 100)
        var range = max - min;

        // Helper converts pixel positions into range position and calls onValueChange
        var updatePosition = function (mouseX, mouseY) {

            var r = el.getBoundingClientRect();

            // Ensure mouse / finger is within control
            if (mouseX < r.left || mouseX > r.right || mouseY > r.bottom || mouseY < r.top) {
                return;
            }

            // Where is the mouse positioned within track
            var pos;
            if ((r.bottom - mouseY) < padding) {
                // Position is in the lower dead zone
                pos = 0;
            } else if (mouseY < (r.bottom - r.height + padding)) {
                // Position is in the upper dead zone
                pos = max;
            } else {
                // Calculate position in pixels
                pos = (r.bottom - padding) - mouseY;
                // Convert to value between min and max
                pos = (pos / (r.height - 2 * padding)) * range + min;
            }
            // Inform the caller
            onValueChange(pos);
        };

        // === Touch support for iOS ===

        // Respond to initial touch
        el.addEventListener('touchstart', function (evt) {
            updatePosition(evt.targetTouches[0].clientX, evt.targetTouches[0].clientY);
            evt.preventDefault();
        });

        // Track touch movement over track
        el.addEventListener('touchmove', function (evt) {
            updatePosition(evt.targetTouches[0].clientX, evt.targetTouches[0].clientY);
            evt.preventDefault();
        });

        // === Mouse support for desktop browsers ===

        if ('ontouchstart' in window)
            return;

        var mousedown = false;

        // Track mouse downs over document
        el.addEventListener('mousedown', function (evt) {
            console.info('Mouse Down')
            mousedown = true;
            updatePosition(evt.clientX, evt.clientY);
        });

        // Track mouse up over document
        el.addEventListener('mouseup', function () {
            console.info('Mouse Up')
            mousedown = false;
        }, true);

        // Track mouse move over slider track
        el.addEventListener('mousemove', function (evt) {
            if (mousedown) {
                console.info('Update')
                updatePosition(evt.clientX, evt.clientY);
            }
        }, true);
    }

    render() {
        let range = this.props.max - this.props.min;
        return <div ref={(domElement) => {
            this.domElement = domElement
        }}
                    className={this.props.className}>
            <div className='vs-slider-track' style={{backgroundColor: this.props.sliderColour}}>
                <div className='vs-slider-thumb-container'
                     style={{top: ((range - this.props.pos) / (range / 100)) + '%'}}>
                    <div className='vs-slider-thumb'/>
                    <div className='vs-slider-thumb-anno'>{this.props.pos.toFixed(0)}</div>
                </div>
            </div>
        </div>
    }
}
