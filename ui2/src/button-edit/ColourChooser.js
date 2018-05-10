import * as React from "react";
import './ColourChooser.css';
import {Modal} from "semantic-ui-react";
import "../colour/ColourConversion.js"
import VerticalSlider from "../buttons/VerticalSlider";
import {asColourObject, asColourString, HSVtoRGB, RGBtoHSV} from "../colour/ColourConversion";
import MountNotifier from "./MountNotifier";
import PropTypes from 'prop-types';

export default class ColourChooser extends React.Component {

    static propTypes = {
        // initial colour as string, number or rgb object
        colour: PropTypes.oneOfType([PropTypes.number, PropTypes.string, PropTypes.object]).isRequired,

        // trigger (button) for Modal
        trigger: PropTypes.element.isRequired,

        // function(colour) to call when ok hit
        onOk: PropTypes.func.isRequired,

        // function to call when cancel hit
        onCancel: PropTypes.func
    };

    constructor(props) {
        super(props);

        this.state = {
            hue: 0,
            saturation: 0,
            value: 0
        };
    }

    setInitialColour() {
        let initialColour = RGBtoHSV(asColourObject(this.props.colour));
        this.setState({hue: initialColour.h, saturation: initialColour.s, value: initialColour.v})
    }

    render() {
        let colour = asColourString(HSVtoRGB(this.state.hue, this.state.saturation, this.state.value));
        let saturatedColour = asColourString(HSVtoRGB(this.state.hue, 100, 100));
        return <Modal trigger={this.props.trigger}
                      header='Choose Colour'
                      content={<div className='scrolling content'>
                          <MountNotifier componentDidMount={() => this.setInitialColour()}/>
                          <div className='cp-content-height'>
                              <ColourSwatch className='cp-colour-swatch' colour={colour}/>
                              <div className='cp-sliders'>
                                  <VerticalSlider className='cp-slider cp-slider-hue'
                                                  min={0}
                                                  max={360}
                                                  pos={this.state.hue}
                                                  onPosChange={p => this.setState({hue: p})}/>

                                  <VerticalSlider className='cp-slider cp-slider-saturation'
                                                  min={0}
                                                  max={100}
                                                  pos={this.state.saturation}
                                                  onPosChange={p => this.setState({saturation: p})}
                                                  sliderColour={saturatedColour}/>

                                  <VerticalSlider className='cp-slider cp-slider-value'
                                                  min={0}
                                                  max={100}
                                                  pos={this.state.value}
                                                  onPosChange={p => this.setState({value: p})}
                                                  sliderColour={saturatedColour}/>
                              </div>
                          </div>
                      </div>}
                      actions={[
                          {key: 'ok', content: 'OK', primary: true},
                          {key: 'cancel', content: 'Cancel'}
                      ]}
                      onActionClick={e => e.target.textContent === 'OK' ?
                          this.props.onOk(colour) :
                          (this.props.onCancel === undefined ? null : this.props.onCancel())}/>
    }
}

function ColourSwatch(props) {
    return <div className='cp-colour-swatch' style={{backgroundColor: props.colour}}/>
}
