import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { User } from '../api/user';
import { CartItem } from '../api/cart';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  setAuth: (user: User, token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      setAuth: (user, token) => {
        localStorage.setItem('token', token);
        set({ user, token, isAuthenticated: true });
      },
      logout: () => {
        localStorage.removeItem('token');
        set({ user: null, token: null, isAuthenticated: false });
      },
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token }),
    }
  )
);

interface CartState {
  items: CartItem[];
  setCart: (items: CartItem[]) => void;
  addItem: (item: CartItem) => void;
  updateItem: (product_id: number, quantity: number) => void;
  removeItem: (product_id: number) => void;
  clearCart: () => void;
  getTotal: () => number;
}

export const useCartStore = create<CartState>((set, get) => ({
  items: [],
  setCart: (items) => set({ items }),
  addItem: (item) =>
    set((state) => {
      const existing = state.items.find((i) => i.product_id === item.product_id);
      if (existing) {
        return {
          items: state.items.map((i) =>
            i.product_id === item.product_id
              ? { ...i, quantity: i.quantity + item.quantity }
              : i
          ),
        };
      }
      return { items: [...state.items, item] };
    }),
  updateItem: (product_id, quantity) =>
    set((state) => ({
      items: state.items.map((i) =>
        i.product_id === product_id ? { ...i, quantity } : i
      ),
    })),
  removeItem: (product_id) =>
    set((state) => ({
      items: state.items.filter((i) => i.product_id !== product_id),
    })),
  clearCart: () => set({ items: [] }),
  getTotal: () =>
    get().items.reduce((sum, item) => sum + item.product_price * item.quantity, 0),
}));
