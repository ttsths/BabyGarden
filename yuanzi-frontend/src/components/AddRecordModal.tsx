/**
 * 添加记录模态框组件
 * 基于 Stitch 设计: yuanzi_add_record_modal
 */

import React, { useState } from 'react';
import { cn } from '../utils/cn';

interface AddRecordModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (record: any) => void;
}

export const AddRecordModal: React.FC<AddRecordModalProps> = ({
  isOpen,
  onClose,
  onSave,
}) => {
  const [selectedType, setSelectedType] = useState<'feeding' | 'sleep' | 'diaper'>('feeding');
  const [time] = useState('10:30 AM');
  const [amount] = useState(120);
  const [note, setNote] = useState('');

  if (!isOpen) return null;

  const handleSave = () => {
    const record = {
      type: selectedType,
      time,
      amount: selectedType === 'feeding' ? amount : undefined,
      note,
    };
    onSave(record);
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/40 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Bottom Sheet Modal */}
      <div className="absolute inset-x-0 bottom-0 z-20 bg-background-light dark:bg-background-dark rounded-t-xl shadow-2xl flex flex-col max-h-[90%] animate-slide-up">
        {/* Sheet Handle */}
        <button
          className="flex h-6 w-full items-center justify-center pt-2"
          onClick={onClose}
        >
          <div className="h-1.5 w-12 rounded-full bg-primary/30" />
        </button>

        <div className="flex-1 overflow-y-auto pb-8">
          {/* Title */}
          <h1 className="text-slate-900 dark:text-slate-100 text-xl font-bold leading-tight px-6 pt-4 pb-6 text-center">
            添加记录
          </h1>

          {/* Record Type Selector */}
          <div className="grid grid-cols-3 gap-4 px-6 mb-8">
            <RecordTypeSelector
              icon="medical_services"
              label="喂奶"
              active={selectedType === 'feeding'}
              onClick={() => setSelectedType('feeding')}
            />
            <RecordTypeSelector
              icon="bedtime"
              label="睡眠"
              active={selectedType === 'sleep'}
              onClick={() => setSelectedType('sleep')}
            />
            <RecordTypeSelector
              icon="baby_changing_station"
              label="换尿布"
              active={selectedType === 'diaper'}
              onClick={() => setSelectedType('diaper')}
            />
          </div>

          {/* Input Fields */}
          <div className="space-y-4 px-6">
            {/* Time Picker */}
            <div className="flex items-center justify-between bg-white dark:bg-slate-800 p-4 rounded-xl shadow-sm border border-primary/5">
              <div className="flex items-center gap-3">
                <span className="material-symbols-outlined text-primary text-xl">schedule</span>
                <p className="text-slate-900 dark:text-slate-100 text-base font-medium">时间</p>
              </div>
              <div className="px-3 py-1 bg-primary/10 rounded-lg">
                <p className="text-primary font-bold">{time}</p>
              </div>
            </div>

            {/* Slider for Amount (only for feeding) */}
            {selectedType === 'feeding' && (
              <div className="bg-white dark:bg-slate-800 p-4 rounded-xl shadow-sm border border-primary/5">
                <div className="flex w-full items-center justify-between mb-4">
                  <div className="flex items-center gap-3">
                    <span className="material-symbols-outlined text-primary text-xl">opacity</span>
                    <p className="text-slate-900 dark:text-slate-100 text-base font-medium">奶量 (ml)</p>
                  </div>
                  <p className="text-primary font-bold text-lg">{amount}</p>
                </div>
                <div className="relative flex h-6 w-full items-center">
                  <div className="h-2 flex-1 rounded-full bg-primary/20">
                    <div
                      className={cn(
                        'h-full rounded-full bg-primary relative transition-all',
                        { 'w-[45%]': amount === 120 },
                        { 'w-[30%]': amount === 80 },
                        { 'w-[60%]': amount === 160 },
                        { 'w-[75%]': amount === 200 },
                      )}
                      style={{ width: `${(amount / 300) * 100}%` }}
                    >
                      <div className="absolute right-0 -top-2 size-6 rounded-full bg-white border-2 border-primary shadow-md" />
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Notes Field */}
            <div className="bg-white dark:bg-slate-800 p-4 rounded-xl shadow-sm border border-primary/5">
              <div className="flex items-center gap-3 mb-3">
                <span className="material-symbols-outlined text-primary text-xl">edit_note</span>
                <p className="text-slate-900 dark:text-slate-100 text-base font-medium">备注</p>
              </div>
              <textarea
                value={note}
                onChange={(e) => setNote(e.target.value)}
                placeholder="写下点什么..."
                className="w-full bg-background-light dark:bg-slate-700 border-none rounded-lg p-3 text-sm focus:ring-1 focus:ring-primary/50 text-slate-700 h-20 resize-none"
                rows={4}
              />
            </div>
          </div>

          {/* Submit Button */}
          <div className="px-6 mt-10">
            <button
              onClick={handleSave}
              className="w-full bg-primary hover:bg-primary/90 text-white font-bold py-4 rounded-xl shadow-lg shadow-primary/30 transition-all flex items-center justify-center gap-2 active:scale-95"
            >
              <span className="material-symbols-outlined">save</span>
              保存记录
            </button>
            {/* Spacer for iOS home indicator */}
            <div className="h-4" />
          </div>
        </div>
      </div>
    </div>
  );
};

// 记录类型选择器
const RecordTypeSelector: React.FC<{
  icon: string;
  label: string;
  active: boolean;
  onClick: () => void;
}> = ({ icon, label, active, onClick }) => (
  <div
    onClick={onClick}
    className="flex flex-col items-center gap-2 group cursor-pointer"
  >
    <div
      className={cn(
        'w-full aspect-square rounded-full flex items-center justify-center shadow-lg transition-all',
        active
          ? 'bg-primary ring-4 ring-primary/10'
          : 'bg-white border-2 border-primary/20'
      )}
    >
      <span className={cn('material-symbols-outlined text-3xl', active ? 'text-white' : 'text-primary')}>
        {icon}
      </span>
    </div>
    <p className={cn('text-sm', active ? 'font-semibold text-slate-900' : 'font-medium text-slate-500')}>
      {label}
    </p>
  </div>
);
