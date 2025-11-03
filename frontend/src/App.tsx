import {CSSProperties, useEffect, useMemo, useState} from 'react';
import './App.css';
import {PetSprite} from "../wailsjs/go/app/App";
import {WindowSetAlwaysOnTop} from "../wailsjs/runtime/runtime";

const dragStyle: CSSProperties = {
    ['--wails-draggable' as any]: 'drag',
};

const noDragStyle: CSSProperties = {
    ['--wails-draggable' as any]: 'no-drag',
};

function App() {
    const [sprite, setSprite] = useState<string>('');
    const [isFloating, setIsFloating] = useState(true);
    const [showHint, setShowHint] = useState(true);

    useEffect(() => {
        let isMounted = true;

        document.body.style.background = "transparent";
        console.log("body background:", getComputedStyle(document.body).backgroundColor);

        WindowSetAlwaysOnTop(true);

        PetSprite()
            .then((data) => {
                if (isMounted) {
                    setSprite(data);
                }
            })
            .catch((err) => console.error('Pet sprite failed to load', err));

        const hintTimeout = window.setTimeout(() => setShowHint(false), 6000);

        return () => {
            isMounted = false;
            window.clearTimeout(hintTimeout);
        };
    }, []);

    const petClassName = useMemo(
        () => `pet-image${isFloating ? ' pet-image--floating' : ''}`,
        [isFloating],
    );

    const toggleFloating = () => setIsFloating((prev) => !prev);

    return (
        <div className="app-shell" style={dragStyle}>
            <div className="pet-container" style={dragStyle}>
                {sprite && (
                    <img
                        className={petClassName}
                        src={sprite}
                        alt="Desktop pet"
                        draggable={false}
                        onDoubleClick={toggleFloating}
                        style={noDragStyle}
                    />
                )}
                <div className="pet-shadow" style={noDragStyle}/>
            </div>

            {showHint && (
                <div className="pet-hint" style={noDragStyle}>
                    Drag me anywhere · Double-click to pause
                </div>
            )}
        </div>
    );
}

export default App;
