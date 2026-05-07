import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/Button';
import { Input } from '@/components/ui/Input';
import { Card } from '@/components/ui/Card';
import { useBabyStore } from '@/stores/useBabyStore';

/**
 * 宝宝信息设置页面
 */
export const BabySetupPage: React.FC = () => {
  const navigate = useNavigate();
  const { addBaby } = useBabyStore();
  const [formData, setFormData] = useState({
    name: '',
    birthday: '',
    gender: 'male' as 'male' | 'female',
  });
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!formData.name || !formData.birthday) return;

    setLoading(true);
    try {
      await addBaby({
        name: formData.name,
        birthday: formData.birthday,
        gender: formData.gender,
      });
      navigate('/');
    } catch (error) {
      console.error('添加宝宝失败:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background-primary flex flex-col">
      {/* 顶部装饰 */}
      <div className="bg-gradient-to-b from-brand-light/20 to-transparent pt-12 pb-8">
        <div className="text-center">
          <div className="w-20 h-20 mx-auto mb-4 rounded-full bg-brand-primary/20 flex items-center justify-center text-4xl">
            👶
          </div>
          <h1 className="text-2xl font-bold text-neutral-text-primary">
            欢迎加入圆子
          </h1>
          <p className="text-neutral-text-secondary mt-2">
            记录宝宝成长的每一个珍贵瞬间
          </p>
        </div>
      </div>

      {/* 表单 */}
      <main className="flex-1 px-6 py-8">
        <Card padding="large">
          <div className="space-y-6">
            {/* 宝宝昵称 */}
            <Input
              label="宝宝昵称"
              placeholder="例如：圆子"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />

            {/* 出生日期 */}
            <Input
              label="出生日期"
              type="date"
              value={formData.birthday}
              onChange={(e) => setFormData({ ...formData, birthday: e.target.value })}
            />

            {/* 性别 */}
            <div>
              <label className="text-neutral-text-primary font-medium mb-3 block">
                性别
              </label>
              <div className="grid grid-cols-2 gap-3">
                <button
                  onClick={() => setFormData({ ...formData, gender: 'male' })}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    formData.gender === 'male'
                      ? 'border-brand-primary bg-brand-light/10'
                      : 'border-neutral-border'
                  }`}
                >
                  <span className="text-2xl">👦</span>
                  <span className="block mt-2 text-neutral-text-primary font-medium">
                    男宝
                  </span>
                </button>
                <button
                  onClick={() => setFormData({ ...formData, gender: 'female' })}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    formData.gender === 'female'
                      ? 'border-brand-primary bg-brand-light/10'
                      : 'border-neutral-border'
                  }`}
                >
                  <span className="text-2xl">👧</span>
                  <span className="block mt-2 text-neutral-text-primary font-medium">
                    女宝
                  </span>
                </button>
              </div>
            </div>

            {/* 提交按钮 */}
            <Button
              variant="primary"
              size="large"
              fullWidth
              loading={loading}
              disabled={!formData.name || !formData.birthday || loading}
              onClick={handleSubmit}
              className="mt-6"
            >
              开始记录
            </Button>
          </div>
        </Card>
      </main>
    </div>
  );
};

export default BabySetupPage;
