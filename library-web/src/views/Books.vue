<template>
  <div>
    <!-- 顶部操作栏 -->
    <div class="card toolbar">
      <div class="flex gap-2 flex-wrap">
        <input v-model="params.title" class="input" placeholder="书名搜索..." @keyup.enter="search" />
        <input v-model="params.author" class="input" placeholder="作者搜索..." @keyup.enter="search" />
        <select v-model="params.category_id" class="input" @change="search">
          <option :value="null">全部分类</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
        <select v-model="params.available_only" class="input" @change="search">
          <option :value="null">库存状态</option>
          <option :value="true">仅显示可借</option>
          <option :value="false">仅显示缺货</option>
        </select>
        <button class="btn" @click="search">🔍 搜索</button>
        <button class="btn btn-secondary" @click="resetFilters">重置</button>
      </div>
      <div class="flex gap-2" v-if="isAdmin">
        <button class="btn" @click="openAddModal">➕ 新增图书</button>
      </div>
    </div>

    <!-- 图书表格 -->
    <div class="card" style="padding: 0; overflow: hidden;">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
      
      <table v-else-if="books.length > 0">
        <thead>
          <tr>
            <th width="80">封面</th>
            <th>书名 / ISBN</th>
            <th>作者 / 出版社</th>
            <th>分类</th>
            <th>价格</th>
            <th>库存状态</th>
            <th>借阅量</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="book in books" :key="book.id" class="book-row">
            <!-- 封面列：点击放大 -->
            <td 
              class="cover-cell" 
              @click="openCoverModal(book)"
              title="点击查看图书封面"
            >
              <div class="cover-wrapper">
                <img 
                  :src="'http://localhost:8080' + (book.cover_url || defaultCover)" 
                  class="book-cover" 
                  @error="handleImgError"
                />
                <div class="cover-hint">🔍</div>
              </div>
            </td>
            
            <!-- 书名列：点击跳转详情 -->
            <td 
              class="clickable-cell" 
              @click="viewDetails(book.id)"
              title="查看图书详情"
            >
              <div class="book-title">{{ book.title }}</div>
              <div class="text-sm text-secondary">{{ book.isbn }}</div>
            </td>
            
            <!-- 作者列：点击跳转详情 -->
            <td 
              class="clickable-cell" 
              @click="viewDetails(book.id)"
              title="查看图书详情"
            >
              <div>{{ book.author }}</div>
              <div class="text-sm text-secondary">{{ book.publisher }}</div>
            </td>
            
            <!-- 分类列：点击跳转详情 -->
            <td 
              class="clickable-cell" 
              @click="viewDetails(book.id)"
              title="查看图书详情"
            >
              <span class="badge badge-secondary">{{ book.category_name || '未分类' }}</span>
            </td>
            
            <!-- 价格列：点击跳转详情 -->
            <td 
              class="clickable-cell" 
              @click="viewDetails(book.id)"
              title="查看图书详情"
            >
              <span v-if="book.price">¥{{ book.price.toFixed(2) }}</span>
              <span v-else class="text-secondary">-</span>
            </td>
            
            <!-- 库存列：点击跳转详情 -->
            <td 
              class="clickable-cell" 
              @click="viewDetails(book.id)"
              title="查看图书详情"
            >
              <span v-if="book.available > 0" class="badge badge-success">
                可借 {{ book.available }} / {{ book.stock }}
              </span>
              <span v-else class="badge badge-danger">暂时缺货</span>
            </td>
            
            <!-- 借阅量列：点击跳转详情 -->
            <td 
              class="clickable-cell" 
              @click="viewDetails(book.id)"
              title="查看图书详情"
            >
              <span class="text-primary">{{ book.borrow_count || 0 }}</span> 次
            </td>
            
            <!-- 操作列：不触发跳转 -->
            <td class="text-right action-cell">
              <div class="flex justify-end gap-2">
                <button class="btn btn-sm" v-if="book.available > 0" @click.stop="handleBorrow(book)">
                  借阅
                </button>
                <button class="btn btn-sm btn-secondary" v-else @click.stop="handleReserve(book.id)">
                  预约
                </button>
                <template v-if="isAdmin">
                  <button class="btn btn-sm btn-secondary" @click.stop="openEditModal(book)">编辑</button>
                  <button class="btn btn-sm btn-danger" @click.stop="handleDelete(book)">删除</button>
                </template>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-else class="empty-state">
        <p>📚 暂无相关图书数据</p>
        <p class="text-sm text-secondary">尝试调整搜索条件或新增图书</p>
      </div>

      <!-- 分页 -->
      <div class="pagination" v-if="books.length > 0">
        <div class="page-info">
          共 <strong>{{ total }}</strong> 本图书，第 {{ params.page }} / {{ totalPages }} 页
        </div>
        <div class="page-controls">
          <button class="btn btn-secondary btn-sm" :disabled="params.page <= 1" @click="changePage(-1)">上一页</button>
          <button class="btn btn-secondary btn-sm" :disabled="params.page >= totalPages" @click="changePage(1)">下一页</button>
        </div>
      </div>
    </div>

    <!-- 封面放大模态框 -->
    <div v-if="showCoverModal" class="modal-overlay" @click="showCoverModal = false">
      <div class="cover-modal" @click.stop>
        <button class="close-btn" @click="showCoverModal = false">✕</button>
        <img 
          :src="'http://localhost:8080' + (currentBook?.cover_url || defaultCover)" 
          :alt="currentBook?.title"
          class="modal-cover"
          @error="handleImgError"
        />
        <div class="modal-info">
          <h3>{{ currentBook?.title }}</h3>
          <p>👤 {{ currentBook?.author }}</p>
        </div>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal" style="width: 550px;">
        <h2>{{ modalType === 'add' ? '📚 图书入库' : '✏️ 编辑图书' }}</h2>
        
        <div class="form-row">
          <div class="form-group">
            <label>书名 <span class="required">*</span></label>
            <input v-model="bookForm.title" class="input" placeholder="1-200字符">
          </div>
          <div class="form-group">
            <label>作者 <span class="required">*</span></label>
            <input v-model="bookForm.author" class="input" placeholder="1-100字符">
          </div>
        </div>
        
        <div class="form-group">
          <label>ISBN <span class="required">*</span></label>
          <input v-model="bookForm.isbn" class="input" placeholder="格式：978-X-XXX-XXXXX-X" :disabled="modalType === 'edit'">
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label>分类 <span class="required">*</span></label>
            <select v-model="bookForm.category_id" class="input">
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
          <div class="form-group">
            <label>出版社 <span class="required">*</span></label>
            <input v-model="bookForm.publisher" class="input">
          </div>
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label>出版日期</label>
            <input v-model="bookForm.publish_date" type="date" class="input">
          </div>
          <div class="form-group">
            <label>价格 (元)</label>
            <input v-model.number="bookForm.price" type="number" step="0.01" class="input" placeholder="0.00">
          </div>
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label>库存数量 <span class="required">*</span></label>
            <input v-model.number="bookForm.stock" type="number" min="0" class="input">
          </div>
          <div class="form-group">
            <label>封面URL</label>
            <input v-model="bookForm.cover_url" class="input" placeholder="https://...">
          </div>
        </div>
        
        <div class="form-group">
          <label>图书简介</label>
          <textarea v-model="bookForm.description" class="input" rows="3" placeholder="最多1000字符" style="height: auto;"></textarea>
        </div>
        
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn" :disabled="submitting" @click="submitBook">
            {{ submitting ? '提交中...' : '确认提交' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { getBooks, borrowBook, reserveBook, deleteBook, addBook, updateBook, getCategories } from '../api';
import { getUser } from '../utils/auth';
import { $message } from '../utils/toast';

const router = useRouter();
const defaultCover = 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" width="150" height="200" viewBox="0 0 150 200"%3E%3Crect fill="%23f3f4f6" width="150" height="200"/%3E%3Ctext fill="%239ca3af" font-family="Arial" font-size="14" x="50%25" y="50%25" text-anchor="middle" dy=".3em"%3E暂无封面%3C/text%3E%3C/svg%3E';

const books = ref([]);
const categories = ref([]);
const total = ref(0);
const totalPages = ref(1);
const loading = ref(false);
const submitting = ref(false);

const params = reactive({ 
  page: 1, 
  limit: 10, 
  title: '', 
  author: '', 
  category_id: null,
  available_only: null
});

const showModal = ref(false);
const modalType = ref('add');
const editingId = ref(null);

const showCoverModal = ref(false);
const currentBook = ref(null);

const bookForm = reactive({
  title: '', author: '', isbn: '', category_id: 1, publisher: '',
  publish_date: '', price: null, stock: 10, description: '', cover_url: ''
});

const user = getUser();
const isAdmin = computed(() => user?.role === 'admin');

// 打开封面放大模态框
const openCoverModal = (book) => {
  currentBook.value = book;
  showCoverModal.value = true;
};

// 查看图书详情
const viewDetails = (bookId) => {
  router.push(`/books/${bookId}`);
};

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await getCategories({ include_count: true });
    categories.value = res.categories || [];
  } catch (e) {
    console.error('加载分类失败:', e);
  }
};

const fetchBooks = async () => {
  loading.value = true;
  try {
    // 清理空参数
    const query = { page: params.page, limit: params.limit };
    if (params.title) query.title = params.title;
    if (params.author) query.author = params.author;
    if (params.category_id) query.category_id = params.category_id;
    if (params.available_only !== null) query.available_only = params.available_only;

    const res = await getBooks(query);
    books.value = res.books || [];
    total.value = res.total || 0;
    totalPages.value = res.total_pages || 1;
  } catch (e) {
    $message.error('数据加载失败');
    books.value = [];
  } finally {
    loading.value = false;
  }
};

const search = () => {
  params.page = 1;
  fetchBooks();
};

const resetFilters = () => {
  Object.assign(params, { page: 1, title: '', author: '', category_id: null, available_only: null });
  fetchBooks();
};

const changePage = (delta) => {
  params.page += delta;
  fetchBooks();
};

const handleImgError = (e) => {
  e.target.src = defaultCover;
};

const openAddModal = () => {
  modalType.value = 'add';
  Object.assign(bookForm, {
    title: '', author: '', isbn: '', category_id: categories.value[0]?.id || 1,
    publisher: '', publish_date: '', price: null, stock: 10, description: '', cover_url: ''
  });
  showModal.value = true;
};

const openEditModal = (book) => {
  modalType.value = 'edit';
  editingId.value = book.id;
  Object.assign(bookForm, {
    title: book.title,
    author: book.author,
    isbn: book.isbn,
    category_id: book.category_id,
    publisher: book.publisher,
    publish_date: book.publish_date || '',
    price: book.price,
    stock: book.stock,
    description: book.description || '',
    cover_url: book.cover_url || ''
  });
  showModal.value = true;
};

const submitBook = async () => {
  // 基础验证
  if (!bookForm.title || !bookForm.author || !bookForm.isbn || !bookForm.publisher) {
    $message.warning('请填写所有必填字段');
    return;
  }

  submitting.value = true;
  try {
    const data = { ...bookForm };
    // 清理空值
    if (!data.price) delete data.price;
    if (!data.publish_date) delete data.publish_date;
    if (!data.description) delete data.description;
    if (!data.cover_url) delete data.cover_url;

    if (modalType.value === 'add') {
      await addBook(data);
      $message.success('图书入库成功');
    } else {
      // 编辑时不传ISBN
      delete data.isbn;
      await updateBook(editingId.value, data);
      $message.success('图书更新成功');
    }
    showModal.value = false;
    fetchBooks();
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
};

const handleBorrow = async (book) => {
  if (confirm(`确认借阅《${book.title}》?\n默认借阅30天`)) {
    try {
      await borrowBook({ book_id: book.id, borrow_days: 30 });
      $message.success('借阅成功！请按时归还');
      fetchBooks();
    } catch (e) {}
  }
};

const handleReserve = async (bookId) => {
  try {
    await reserveBook({ book_id: bookId });
    $message.success('预约成功！书归还后将通知您');
  } catch (e) {}
};

const handleDelete = async (book) => {
  if (confirm(`确定删除《${book.title}》?\n此操作不可恢复。`)) {
    try {
      await deleteBook(book.id);
      $message.success('图书删除成功');
      fetchBooks();
    } catch (e) {}
  }
};

onMounted(() => {
  loadCategories();
  fetchBooks();
});
</script>

<style scoped>
.toolbar { 
  display: flex; 
  justify-content: space-between; 
  align-items: center;
  padding: 16px;
  flex-wrap: wrap;
  gap: 12px;
}
.toolbar .input { width: 160px; }

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--border-light);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.book-row {
  transition: var(--transition);
}

/* 封面单元格样式 */
.cover-cell {
  cursor: pointer;
  padding: 12px !important;
}

.cover-wrapper {
  position: relative;
  width: 48px;
  height: 64px;
  border-radius: 4px;
  overflow: hidden;
}

.book-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: var(--transition);
}

.cover-hint {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  opacity: 0;
  transition: opacity 0.3s;
}

.cover-cell:hover .cover-hint {
  opacity: 1;
}

.cover-cell:hover .book-cover {
  transform: scale(1.1);
}

/* 可点击单元格样式 */
.clickable-cell {
  cursor: pointer;
  position: relative;
  transition: var(--transition);
}

.clickable-cell:hover {
  background: var(--border-light);
}

/* 操作单元格不显示提示 */
.action-cell {
  cursor: default;
}

.book-title { 
  font-weight: 600; 
  max-width: 200px; 
  overflow: hidden; 
  text-overflow: ellipsis; 
  white-space: nowrap; 
}

.text-sm { 
  font-size: 0.8rem; 
}

.text-secondary { 
  color: var(--text-secondary); 
}

.text-primary { 
  color: var(--primary); 
  font-weight: 600; 
}

.empty-state { 
  padding: 60px; 
  text-align: center; 
  color: var(--text-secondary); 
}

.empty-state p { 
  margin: 8px 0; 
}

.pagination { 
  padding: 16px; 
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-info { 
  color: var(--text-secondary); 
  font-size: 0.875rem; 
}

.page-controls { 
  display: flex; 
  gap: 8px; 
}

/* 封面放大模态框 */
.cover-modal {
  background: white;
  border-radius: var(--radius-lg);
  padding: 32px;
  max-width: 600px;
  width: 90%;
  position: relative;
  animation: modalSlideIn 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 36px;
  height: 36px;
  border: none;
  background: rgba(0, 0, 0, 0.1);
  color: var(--text-main);
  border-radius: 50%;
  font-size: 1.2rem;
  cursor: pointer;
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.2);
  transform: scale(1.1);
}

.modal-cover {
  width: 100%;
  max-height: 600px;
  object-fit: contain;
  border-radius: var(--radius);
  margin-bottom: 20px;
}

.modal-info h3 {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 8px 0;
}

.modal-info p {
  font-size: 1rem;
  color: var(--text-secondary);
  margin: 0;
}

.form-row { 
  display: grid; 
  grid-template-columns: 1fr 1fr; 
  gap: 16px; 
}

.required { 
  color: var(--danger); 
}

.flex-wrap { 
  flex-wrap: wrap; 
}
</style>
