import React, { useState } from 'react';
import { Button } from '../components/Button';
import { Card } from '../components/Card';
import { RECORD_TYPES } from '../types';
import type { RecordType } from '../types';

interface RecordPageProps {
  onSave: (record: Record<string, unknown>) => void;
  onCancel: () => void;
}

export const RecordPage: React.FC<RecordPageProps> = ({ onSave, onCancel }) => {
  const [selectedType, setSelectedType] = useState<RecordType | null>(null);
  const [amount, setAmount] = useState<number>(120);
  const [side, setSide] = useState<'left' | 'right' | 'both'>('both');
  const [note, setNote] = useState('');

  const handleTypeSelect = (type: RecordType) => {
    setSelectedType(type);
  };

  const handleSave = () => {
    if (!selectedType) return;

    const record = {
      type: selectedType.id,
      startedAt: new Date().toISOString(),
      content: selectedType.id === 'feeding' ? {
        type: 'breast',
        side,
        duration: 15,
      } : {},
      note,
    };

    onSave(record);
  };

  return (
    <div className="min-h-screen bg-bg-primary">
      {/* 顶部导航 */}
      <div className="sticky top-0 z-10 bg-bg-primary/95 backdrop-blur-sm border-b border-border">
        <div className="p-lg">
          <div className="flex items-center justify-between">
            <Button variant="ghost" size="small" onClick={onCancel}>
              ← 取消
            </Button>
            <h1 className="text-h1 font-semibold text-text-primary">记录</h1>
            <Button variant="primary" size="small" onClick={handleSave} disabled={!selectedType}>
              保存
            </Button>
          </div>
        </div>
      </div>

      {/* 主要内容 */}
      <main className="p-lg">
        {/* 记录类型选择 */}
        {!selectedType && (
          <div>
            <p className="text-h3 font-medium text-text-primary mb-md">选择记录类型</p>
            <div className="grid grid-cols-2 gap-sm">
              {RECORD_TYPES.map((type) => (
                <TypeCard
                  key={type.id}
                  type={type}
                  onClick={() => handleTypeSelect(type)}
                />
              ))}
            </div>
          </div>
        )}

        {/* 记录详情表单 */}
        {selectedType && (
          <div className="space-y-lg">
            {/* 顶部返回按钮 */}
            <Button variant="ghost" size="small" onClick={() => setSelectedType(null)}>
              ← 返回类型选择
            </Button>

            {/* 喂养记录表单 */}
            {selectedType.id === 'feeding' && (
              <Card>
                <div className="space-y-lg">
                  <div>
                    <p className="text-h3 font-medium text-text-primary mb-sm">喂养方式</p>
                    <div className="grid grid-cols-2 gap-sm">
                      <FeedTypeButton
                        label="母乳喂养"
                        type="breast"
                        active
                        onClick={() => {}}
                      />
                      <FeedTypeButton
                        label="奶粉喂养"
                        type="formula"
                        onClick={() => {}}
                      />
                    </div>
                  </div>

                  <div>
                    <p className="text-h3 font-medium text-text-primary mb-sm">奶量</p>
                    <div className="flex items-center gap-md">
                      <Button variant="secondary" size="small" onClick={() => setAmount(Math.max(0, amount - 10))}>
                        -
                      </Button>
                      <div className="flex-1">
                        <input
                          type="number"
                          value={amount}
                          onChange={(e) => setAmount(Number(e.target.value))}
                          className="w-full text-center text-h1 font-semibold bg-bg-tertiary rounded-md p-md border border-border"
                        />
                      </div>
                      <span className="text-body text-text-secondary">ml</span>
                      <Button variant="secondary" size="small" onClick={() => setAmount(amount + 10)}>
                        +
                      </Button>
                    </div>
                  </div>

                  <div>
                    <p className="text-h3 font-medium text-text-primary mb-sm">侧别</p>
                    <div className="flex gap-sm">
                      <SideButton label="左" active={side === 'left'} onClick={() => setSide('left')} />
                      <SideButton label="双" active={side === 'both'} onClick={() => setSide('both')} />
                      <SideButton label="右" active={side === 'right'} onClick={() => setSide('right')} />
                    </div>
                  </div>
                </div>
              </Card>
            )}

            {/* 睡眠记录表单 */}
            {selectedType.id === 'sleep' && (
              <Card>
                <div className="space-y-lg">
                  <div className="text-center">
                    <div className="text-6xl mb-md">💤</div>
                    <p className="text-h1 font-semibold text-text-primary mb-sm">睡眠记录</p>
                    <p className="text-body text-text-secondary">记录宝宝的睡眠时间</p>
                  </div>
                </div>
              </Card>
            )}

            {/* 换尿布记录表单 */}
            {selectedType.id === 'diaper' && (
              <Card>
                <div className="space-y-lg">
                  <div className="text-center">
                    <div className="text-6xl mb-md">💩</div>
                    <p className="text-h1 font-semibold text-text-primary mb-sm">换尿布记录</p>
                    <p className="text-body text-text-secondary">记录宝宝的排泄情况</p>
                  </div>
                </div>
              </Card>
            )}

            {/* 备注输入 */}
            <Card>
              <div className="space-y-sm">
                <p className="text-h3 font-medium text-text-primary">备注（可选）</p>
                <textarea
                  value={note}
                  onChange={(e) => setNote(e.target.value)}
                  placeholder="添加备注..."
                  className="w-full min-h-24 bg-bg-tertiary rounded-md p-md border border-border text-body text-text-primary placeholder:text-text-disabled"
                  rows={4}
                />
              </div>
            </Card>
          </div>
        )}
      </main>
    </div>
  );
};

// 记录类型卡片
const TypeCard: React.FC<{
  type: RecordType;
  onClick: () => void;
}> = ({ type, onClick }) => (
  <Card padding="lg" onClick={onClick} hoverable>
    <div className="text-center space-y-sm">
      <div className="text-5xl">{type.icon}</div>
      <div className="text-body font-medium text-text-primary">{type.name}</div>
    </div>
  </Card>
);

// 喂养方式按钮
const FeedTypeButton: React.FC<{
  label: string;
  type: 'breast' | 'formula';
  active?: boolean;
  onClick: () => void;
}> = ({ label, active, onClick }) => (
  <button
    onClick={onClick}
    className={`
      p-md rounded-lg border-2 transition-all duration-150
      ${active ? 'border-brand bg-brand/10' : 'border-border hover:border-brand/50'}
    `}
  >
    <span className="text-body font-medium text-text-primary">{label}</span>
  </button>
);

// 侧别按钮
const SideButton: React.FC<{
  label: string;
  active?: boolean;
  onClick: () => void;
}> = ({ label, active, onClick }) => (
  <button
    onClick={onClick}
    className={`
      flex-1 p-sm rounded-lg transition-all duration-150
      ${active ? 'bg-brand text-white' : 'bg-bg-tertiary text-text-primary'}
    `}
  >
    <span className="text-body font-medium">{label}</span>
  </button>
);
