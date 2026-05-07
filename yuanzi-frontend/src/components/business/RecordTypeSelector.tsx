import React from 'react';
import { Card } from '@/components/ui/Card';

interface RecordTypeSelectorProps {
  selectedType?: string;
  onSelectType: (typeId: string) => void;
  className?: string;
}

/**
 * 记录类型选择器
 * 用于快速选择记录类型
 */
export const RecordTypeSelector: React.FC<RecordTypeSelectorProps> = ({
  selectedType,
  onSelectType,
  className = '',
}) => {
  const recordTypes = [
    { id: 'feeding', icon: '🍼', label: '喂奶' },
    { id: 'sleep', icon: '💤', label: '睡觉' },
    { id: 'diaper', icon: '💩', label: '换尿布' },
    { id: 'temperature', icon: '🌡️', label: '体温' },
    { id: 'food', icon: '🍚', label: '辅食' },
    { id: 'medicine', icon: '💊', label: '用药' },
    { id: 'bath', icon: '🛁', label: '洗澡' },
    { id: 'other', icon: '📷', label: '其他' },
  ];

  return (
    <Card padding="medium" className={className}>
      <div className="grid grid-cols-4 gap-4">
        {recordTypes.map((type) => (
          <button
            key={type.id}
            onClick={() => onSelectType(type.id)}
            className={`flex flex-col items-center gap-2 p-4 rounded-xl transition-all ${
              selectedType === type.id
                ? 'bg-brand-primary text-white shadow-md'
                : 'bg-background-tertiary text-neutral-text-primary hover:bg-brand-light/20'
            }`}
          >
            <span className="text-3xl">{type.icon}</span>
            <span className="text-sm font-medium">{type.label}</span>
          </button>
        ))}
      </div>
    </Card>
  );
};

export default RecordTypeSelector;
