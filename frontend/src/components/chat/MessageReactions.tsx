'use client';

interface MessageReactionsProps {
    reactions: ReactionCount[];
    currentUserId: string;
    onToggleReaction: (emoji: string) => void;
}

interface ReactionCount {
    emoji: string;
    count: number;
    users: string[];
}

export default function MessageReactions({
    reactions,
    currentUserId,
    onToggleReaction,
}: MessageReactionsProps) {
    if (!reactions || reactions.length === 0) return null;

    return (
        <div className="message-reactions">
            {reactions.map((reaction) => {
                const hasReacted = reaction.users.includes(currentUserId);
                return (
                    <button
                        key={reaction.emoji}
                        onClick={() => onToggleReaction(reaction.emoji)}
                        className={`reaction-badge ${hasReacted ? 'active' : ''}`}
                        title={`${reaction.count} reaction${reaction.count > 1 ? 's' : ''}`}
                    >
                        <span className="emoji">{reaction.emoji}</span>
                        <span className="count">{reaction.count}</span>
                    </button>
                );
            })}

            <style jsx>{`
        .message-reactions {
          display: flex;
          flex-wrap: wrap;
          gap: 4px;
          margin-top: 4px;
        }
        .reaction-badge {
          display: inline-flex;
          align-items: center;
          gap: 4px;
          padding: 2px 8px;
          border-radius: 12px;
          border: 1px solid var(--color-earth-200, #e8dcc8);
          background: var(--color-earth-50, #faf8f5);
          cursor: pointer;
          font-size: 0.85rem;
          transition: all 0.15s;
        }
        .reaction-badge:hover {
          border-color: var(--color-gold-400, #f5c542);
          background: var(--color-gold-50, #fffef5);
        }
        .reaction-badge.active {
          border-color: var(--color-gold-400, #f5c542);
          background: var(--color-gold-100, #fff8e5);
        }
        .emoji {
          font-size: 0.9rem;
        }
        .count {
          color: var(--color-earth-600, #8b7355);
          font-weight: 500;
        }
      `}</style>
        </div>
    );
}
