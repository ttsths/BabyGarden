import { create } from 'zustand';
import type { Baby } from '@/types/models';
import { api } from '@/services/api';

interface BabyState {
  // 状态
  babies: Baby[];
  currentBabyId: string | null;
  isLoading: boolean;
  
  // 动作
  fetchBabies: () => Promise<void>;
  selectBaby: (babyId: string) => void;
  addBaby: (baby: any) => Promise<void>;
  updateBaby: (babyId: string, baby: any) => Promise<void>;
  getCurrentBaby: () => Baby | null;
}

export const useBabyStore = create<BabyState>()((set, get) => ({
  babies: [],
  currentBabyId: null,
  isLoading: false,

  fetchBabies: async () => {
    set({ isLoading: true });
    try {
      const res = await api.baby.getList();
      const babies = res.data as Baby[];
      set({ babies, isLoading: false });
      
      // 自动选择第一个宝宝
      if (babies.length > 0 && !get().currentBabyId) {
        set({ currentBabyId: babies[0].id });
      }
    } catch (error) {
      set({ isLoading: false });
      throw error;
    }
  },

  selectBaby: (babyId) => {
    set({ currentBabyId: babyId });
  },

  addBaby: async (baby) => {
    const newRes = await api.baby.create(baby);
    const newBaby = newRes.data as Baby;
    set((state) => ({
      babies: [...state.babies, newBaby],
      currentBabyId: newBaby.id,
    }));
  },

  updateBaby: async (babyId, baby) => {
    const updatedRes = await api.baby.update(babyId, baby);
    const updatedBaby = updatedRes.data as Baby;
    set((state) => ({
      babies: state.babies.map((b) =>
        b.id === babyId ? updatedBaby : b
      ),
    }));
  },

  getCurrentBaby: () => {
    const { babies, currentBabyId } = get();
    return babies.find((b) => b.id === currentBabyId) || null;
  },
}));
