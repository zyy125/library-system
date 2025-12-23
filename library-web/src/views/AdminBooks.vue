<template>
  <div>
    <div class="card toolbar">
      <div class="flex gap-2 flex-wrap">
        <input v-model="params.title" class="input" placeholder="书名..." @keyup.enter="search" />
        <input v-model="params.author" class="input" placeholder="作者..." @keyup.enter="search" />
        <input v-model="params.isbn" class="input" placeholder="ISBN..." @keyup.enter="search" />
        <select v-model="params.category_id" class="input" @change="search">
          <option :value="null">全部分类</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">
            {{ cat.name }}
          </option>
        </select>
        <button class="btn" @click="search">🔍 搜索</button>
        <button class="btn btn-secondary" @click="resetFilters">重置</button>
      </div>
      <div class="flex gap-2">
        <button class="btn" @click="openModal('create')">➕ 新增图书</button>
        <button class="btn btn-secondary" @click="showBatchModal = true">📥 批量导入</button>
      </div>
    </div>

    <div class="card" style="padding: 0;">
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <table v-else-if="books.length > 0">
        <thead>
          <tr>
            <th width="80">封面</th>
            <th>书名 / ISBN</th>
            <th>作者</th>
            <th>分类</th>
            <th>价格</th>
            <th>库存</th>
            <th>借阅量</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="book in books" :key="book.id">
            <td>
              <img 
                v-if="book.cover_url" 
                :src="book.cover_url" 
                class="book-cover"
                @error="handleImgError"
              />
              <div v-else class="book-cover no-cover">📖</div>
            </td>
            <td>
              <div class="book-title">{{ book.title }}</div>
              <div class="text-sm text-secondary">{{ book.isbn }}</div>
            </td>
            <td>{{ book.author }}</td>
            <td>
              <span class="badge badge-secondary">{{ book.category_name || '未分类' }}</span>
            </td>
            <td>
              <span v-if="book.price">¥{{ book.price.toFixed(2) }}</span>
              <span v-else class="text-secondary">-</span>
            </td>
            <td>
              <div>总量: {{ book.stock }}</div>
              <div class="text-sm" :class="book.available > 0 ? 'text-success' : 'text-danger'">
                可用: {{ book.available }}
              </div>
            </td>
            <td>{{ book.borrow_count || 0 }}</td>
            <td class="text-right">
              <div class="flex justify-end gap-2">
                <button class="btn btn-sm btn-secondary" @click="openModal('edit', book)">
                  编辑
                </button>
                <button class="btn btn-sm btn-danger" @click="handleDelete(book)">
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else class="empty-state">暂无图书数据</div>

      <!-- 分页 -->
      <div class="pagination" v-if="total > 0">
        <div class="page-info">
          共 {{ total }} 本图书，第 {{ params.page }} / {{ totalPages }} 页
        </div>
        <div class="page-controls">
          <button class="btn btn-secondary btn-sm" :disabled="params.page <= 1" @click="changePage(-1)">
            上一页
          </button>
          <button class="btn btn-secondary btn-sm" :disabled="params.page >= totalPages" @click="changePage(1)">
            下一页
          </button>
        </div>
      </div>
    </div>

    <!-- 新增/编辑弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal" style="max-width: 600px;">
        <h3>{{ modalType === 'create' ? '➕ 新增图书' : '✏️ 编辑图书' }}</h3>
        
        <div class="form-row">
          <div class="form-group">
            <label>书名 <span class="required">*</span></label>
            <input v-model="form.title" class="input">
          </div>
          <div class="form-group">
            <label>作者 <span class="required">*</span></label>
            <input v-model="form.author" class="input">
          </div>
        </div>

        <div class="form-group">
          <label>ISBN <span class="required">*</span></label>
          <input 
            v-model="form.isbn" 
            class="input" 
            placeholder="978-X-XXX-XXXXX-X"
            :disabled="modalType === 'edit'"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>分类 <span class="required">*</span></label>
            <select v-model="form.category_id" class="input">
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                {{ cat.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>出版社 <span class="required">*</span></label>
            <input v-model="form.publisher" class="input">
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>出版日期</label>
            <input v-model="form.publish_date" type="date" class="input">
          </div>
          <div class="form-group">
            <label>价格 (元)</label>
            <input v-model.number="form.price" type="number" step="0.01" class="input">
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>库存数量 <span class="required">*</span></label>
            <input v-model.number="form.stock" type="number" min="0" class="input">
          </div>
          <div class="form-group">
            <label>封面URL</label>
            <input v-model="form.cover_url" class="input">
          </div>
        </div>

        <div class="form-group">
          <label>图书简介</label>
          <textarea v-model="form.description" class="input" rows="3" style="height: auto;"></textarea>
        </div>

        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn" :disabled="submitting" @click="submitForm">
            {{ submitting ? '保存中...' : '确认保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 批量导入弹窗 -->
    <div v-if="showBatchModal" class="modal-overlay" @click.self="showBatchModal = false">
      <div class="modal" style="max-width: 700px;">
        <h3>📥 批量导入图书</h3>
        <p class="hint">每行一本书，字段用逗号分隔：书名,作者,ISBN,分类ID,出版社,出版日期,价格,库存,简介</p>
        
        <div class="form-group">
          <label>导入数据</label>
          <textarea 
            v-model="batchText" 
            class="input" 
            rows="10" 
            style="height: auto; font-family: monospace;"
            placeholder="示例：
深入理解计算机系统,Randal E. Bryant,978-7-111-54493-7,1,机械工业出版社,2016-11-01,139.00,10,经典教材
算法导论,Thomas H. Cormen,978-7-111-40701-0,1,机械工业出版社,2012-12-01,128.00,5,算法圣经"
          ></textarea>
        </div>

        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showBatchModal = false">取消</button>
          <button class="btn" :disabled="batchImporting" @click="handleBatchImport">
            {{ batchImporting ? '导入中...' : '开始导入' }}
          </button>
        </div>

        <!-- 导入结果 -->
        <div v-if="batchResult" class="batch-result">
          <h4>导入结果</h4>
          <p>成功: {{ batchResult.success_count }} 本</p>
          <p v-if="batchResult.failed_count > 0" class="text-danger">
            失败: {{ batchResult.failed_count }} 本
          </p>
          <div v-if="batchResult.failed_items?.length > 0" class="failed-list">
            <p><strong>失败项目：</strong></p>
            <ul>
              <li v-for="item in batchResult.failed_items" :key="item.index">
                第 {{ item.index + 1 }} 行 (ISBN: {{ item.isbn }}): {{ item.error }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getBooks, addBook, updateBook, deleteBook, batchImportBooks, getCategories } from '../api';
import { $message } from '../utils/toast';

const books = ref([]);
const categories = ref([]);
const loading = ref(false);
const submitting = ref(false);
const total = ref(0);
const totalPages = ref(1);

const params = reactive({
  page: 1,
  limit: 10,
  title: '',
  author: '',
  isbn: '',
  category_id: null
});

const showModal = ref(false);
const modalType = ref('create');
const currentId = ref(null);

const form = reactive({
  title: '',
  author: '',
  isbn: '',
  category_id: 1,
  publisher: '',
  publish_date: '',
  price: null,
  stock: 10,
  description: '',
  cover_url: ''
});

const showBatchModal = ref(false);
const batchText = ref('');
const batchImporting = ref(false);
const batchResult = ref(null);

const loadCategories = async () => {
  try {
    const res = await getCategories({ include_count: true });
    categories.value = res.categories || [];
  } catch (error) {
    console.error('加载分类失败:', error);
  }
};

const fetchBooks = async () => {
  loading.value = true;
  try {
    const query = { page: params.page, limit: params.limit };
    if (params.title) query.title = params.title;
    if (params.author) query.author = params.author;
    if (params.isbn) query.isbn = params.isbn;
    if (params.category_id) query.category_id = params.category_id;

    const res = await getBooks(query);
    books.value = res.books || [];
    total.value = res.total || 0;
    totalPages.value = res.total_pages || 1;
  } catch (error) {
    console.error('加载图书失败:', error);
  } finally {
    loading.value = false;
  }
};

const search = () => {
  params.page = 1;
  fetchBooks();
};

const resetFilters = () => {
  Object.assign(params, { page: 1, title: '', author: '', isbn: '', category_id: null });
  fetchBooks();
};

const changePage = (delta) => {
  params.page += delta;
  fetchBooks();
};

const handleImgError = (e) => {
  e.target.style.display = 'none';
  e.target.nextElementSibling?.classList.remove('hidden');
};

const openModal = (type, book = null) => {
  modalType.value = type;
  showModal.value = true;
  
  if (type === 'create') {
    Object.assign(form, {
      title: '',
      author: '',
      isbn: '',
      category_id: categories.value[0]?.id || 1,
      publisher: '',
      publish_date: '',
      price: null,
      stock: 10,
      description: '',
      cover_url: ''
    });
  } else {
    currentId.value = book.id;
    Object.assign(form, {
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
  }
};

const submitForm = async () => {
  if (!form.title || !form.author || !form.isbn || !form.publisher) {
    $message.warning('请填写所有必填字段');
    return;
  }

  submitting.value = true;
  try {
    const data = { ...form };
    if (!data.price) delete data.price;
    if (!data.publish_date) delete data.publish_date;
    if (!data.description) delete data.description;
    if (!data.cover_url) delete data.cover_url;

    if (modalType.value === 'create') {
      await addBook(data);
      $message.success('图书添加成功');
    } else {
      delete data.isbn;
      await updateBook(currentId.value, data);
      $message.success('图书更新成功');
    }
    
    showModal.value = false;
    fetchBooks();
  } catch (error) {
    console.error('保存失败:', error);
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (book) => {
  if (!confirm(`确定删除《${book.title}》？此操作不可恢复。`)) return;

  try {
    await deleteBook(book.id);
    $message.success('图书删除成功');
    fetchBooks();
  } catch (error) {
    console.error('删除失败:', error);
  }
};

const handleBatchImport = async () => {
  if (!batchText.value.trim()) {
    $message.warning('请输入导入数据');
    return;
  }

  const lines = batchText.value.trim().split('\n');
  const booksData = [];

  for (const line of lines) {
    const parts = line.split(',').map(p => p.trim());
    if (parts.length < 8) continue;

    booksData.push({
      title: parts[0],
      author: parts[1],
      isbn: parts[2],
      category_id: parseInt(parts[3]),
      publisher: parts[4],
      publish_date: parts[5] || undefined,
      price: parseFloat(parts[6]) || undefined,
      stock: parseInt(parts[7]),
      description: parts[8] || undefined
    });
  }

  if (booksData.length === 0) {
    $message.warning('没有有效的导入数据');
    return;
  }

  batchImporting.value = true;
  batchResult.value = null;
  
  try {
    const res = await batchImportBooks({ books: booksData });
    batchResult.value = res;
    
    if (res.failed_count === 0) {
      $message.success(`成功导入 ${res.success_count} 本图书`);
      setTimeout(() => {
        showBatchModal.value = false;
        batchText.value = '';
        batchResult.value = null;
        fetchBooks();
      }, 2000);
    } else {
      $message.warning(`导入完成：成功 ${res.success_count} 本，失败 ${res.failed_count} 本`);
    }
  } catch (error) {
    console.error('批量导入失败:', error);
  } finally {
    batchImporting.value = false;
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
  gap: 12px;
  flex-wrap: wrap;
}

.toolbar .input {
  width: 140px;
}

.book-cover {
  width: 48px;
  height: 64px;
  object-fit: cover;
  border-radius: 4px;
}

.book-cover.no-cover {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
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

.text-success {
  color: var(--success);
}

.text-danger {
  color: var(--danger);
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

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.required {
  color: var(--danger);
}

.hint {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 12px;
}

.batch-result {
  margin-top: 24px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.batch-result h4 {
  margin-top: 0;
}

.failed-list {
  margin-top: 12px;
  max-height: 200px;
  overflow-y: auto;
}

.failed-list ul {
  margin: 8px 0;
  padding-left: 20px;
}

.failed-list li {
  font-size: 0.875rem;
  margin-bottom: 4px;
}
</style>
