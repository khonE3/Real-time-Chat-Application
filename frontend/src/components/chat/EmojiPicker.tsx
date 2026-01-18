'use client';

import { useState, useRef, useEffect } from 'react';

interface EmojiPickerProps {
    onSelect: (emoji: string) => void;
    onClose: () => void;
}

const EMOJI_CATEGORIES = {
    'Smileys': ['😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂', '🙂', '😊', '😇', '🥰', '😍', '🤩', '😘', '😗', '😚', '😙', '🥲', '😋'],
    'Gestures': ['👍', '👎', '👌', '✌️', '🤞', '🤟', '🤘', '🤙', '👋', '🖐️', '✋', '🖖', '👏', '🙌', '👐', '🤲', '🙏', '💪', '🦾', '🖕'],
    'Hearts': ['❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '🤍', '🤎', '💔', '❣️', '💕', '💞', '💓', '💗', '💖', '💘', '💝', '💟', '♥️'],
    'Reactions': ['👍', '❤️', '😂', '😮', '😢', '😡', '🎉', '🤔', '🙏', '🔥', '💯', '✨', '👀', '🤯', '😱', '🥳', '😎', '🤗', '🥺', '😤'],
    'Objects': ['⭐', '🌟', '💫', '✨', '⚡', '🔥', '💥', '❄️', '🌈', '☀️', '🌙', '⭐', '🎉', '🎊', '🎁', '🎀', '🏆', '🥇', '🏅', '🎯'],
    'Food': ['🍕', '🍔', '🍟', '🌭', '🍿', '🧀', '🥚', '🍳', '🥞', '🧇', '🥓', '🍖', '🍗', '🥩', '🍣', '🍱', '🍜', '🍝', '🍲', '🍛'],
};

const COMMON_REACTIONS = ['👍', '❤️', '😂', '😮', '😢', '😡', '🎉', '🤔'];

export default function EmojiPicker({ onSelect, onClose }: EmojiPickerProps) {
    const [activeCategory, setActiveCategory] = useState('Reactions');
    const [searchQuery, setSearchQuery] = useState('');
    const pickerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (pickerRef.current && !pickerRef.current.contains(e.target as Node)) {
                onClose();
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, [onClose]);

    const handleEmojiClick = (emoji: string) => {
        onSelect(emoji);
        onClose();
    };

    return (
        <div ref={pickerRef} className="emoji-picker">
            {/* Quick reactions */}
            <div className="quick-reactions">
                {COMMON_REACTIONS.map((emoji) => (
                    <button
                        key={emoji}
                        onClick={() => handleEmojiClick(emoji)}
                        className="emoji-btn quick"
                    >
                        {emoji}
                    </button>
                ))}
            </div>

            <div className="divider" />

            {/* Category tabs */}
            <div className="categories">
                {Object.keys(EMOJI_CATEGORIES).map((category) => (
                    <button
                        key={category}
                        onClick={() => setActiveCategory(category)}
                        className={`category-btn ${activeCategory === category ? 'active' : ''}`}
                    >
                        {EMOJI_CATEGORIES[category as keyof typeof EMOJI_CATEGORIES][0]}
                    </button>
                ))}
            </div>

            {/* Emoji grid */}
            <div className="emoji-grid">
                {EMOJI_CATEGORIES[activeCategory as keyof typeof EMOJI_CATEGORIES].map((emoji) => (
                    <button
                        key={emoji}
                        onClick={() => handleEmojiClick(emoji)}
                        className="emoji-btn"
                    >
                        {emoji}
                    </button>
                ))}
            </div>

            <style jsx>{`
        .emoji-picker {
          position: absolute;
          bottom: 100%;
          left: 0;
          margin-bottom: 8px;
          background: white;
          border-radius: 12px;
          box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
          padding: 8px;
          width: 280px;
          z-index: 1000;
        }
        .quick-reactions {
          display: flex;
          gap: 4px;
          padding: 4px;
        }
        .emoji-btn {
          width: 32px;
          height: 32px;
          border: none;
          background: transparent;
          border-radius: 8px;
          cursor: pointer;
          font-size: 1.2rem;
          transition: all 0.15s;
          display: flex;
          align-items: center;
          justify-content: center;
        }
        .emoji-btn:hover {
          background: var(--color-gold-100, #fff8e5);
          transform: scale(1.15);
        }
        .emoji-btn.quick {
          width: 30px;
          height: 30px;
          font-size: 1.3rem;
        }
        .divider {
          height: 1px;
          background: var(--color-earth-200, #e8dcc8);
          margin: 8px 0;
        }
        .categories {
          display: flex;
          gap: 2px;
          padding: 4px;
          overflow-x: auto;
        }
        .category-btn {
          padding: 6px 8px;
          border: none;
          background: transparent;
          border-radius: 8px;
          cursor: pointer;
          font-size: 1rem;
          opacity: 0.6;
          transition: all 0.15s;
        }
        .category-btn:hover,
        .category-btn.active {
          opacity: 1;
          background: var(--color-gold-100, #fff8e5);
        }
        .emoji-grid {
          display: grid;
          grid-template-columns: repeat(8, 1fr);
          gap: 2px;
          max-height: 200px;
          overflow-y: auto;
          padding: 4px;
        }
        .emoji-grid::-webkit-scrollbar {
          width: 4px;
        }
        .emoji-grid::-webkit-scrollbar-thumb {
          background: var(--color-earth-300, #d4c4a8);
          border-radius: 2px;
        }
      `}</style>
        </div>
    );
}
