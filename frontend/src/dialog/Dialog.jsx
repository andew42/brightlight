import * as React from 'react';
import './Dialog.css';

export default function Dialog({open, header, children, actions, onClose}) {

    const ref = React.useRef();

    React.useEffect(() => {
        const dialog = ref.current;
        if (open && !dialog.open) dialog.showModal();
        else if (!open && dialog.open) dialog.close();
    }, [open]);

    // Close on backdrop click — when the click target is the <dialog> element
    // itself (not the inner content wrapper), it must be on the backdrop.
    const handleClick = (e) => {
        if (e.target === ref.current) ref.current.close();
    };

    return (
        <dialog ref={ref} className='dialog' onClose={onClose} onClick={handleClick}>
            <div className='dialog-content'>
                {header && <div className='dialog-header'>{header}</div>}
                <div className='dialog-body'>{children}</div>
                {actions && (
                    <div className='dialog-actions'>
                        {actions.map(a => (
                            <button key={a.key}
                                    className={'dialog-btn' + (a.primary ? ' primary' : '')}
                                    onClick={a.onClick}>
                                {a.content}
                            </button>
                        ))}
                    </div>
                )}
            </div>
        </dialog>
    );
}
