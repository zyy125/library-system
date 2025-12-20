<template>
  <div>
    <!-- 顶部操作栏 -->
    <div class="card flex justify-between" style="padding: 16px;">
      <div class="flex gap-2">
        <input v-model="params.title" class="input" placeholder="输入书名搜索..." @keyup.enter="fetchBooks" />
        <input v-model="params.author" class="input" placeholder="输入作者搜索..." @keyup.enter="fetchBooks" />
        <button class="btn" @click="fetchBooks">
          🔍 搜索
        </button>
      </div>
      <button v-if="isAdmin" class="btn" @click="showAddModal = true">
        ➕ 新增图书
      </button>
    </div>

    <!-- 图书表格 -->
    <div class="card" style="padding: 0; overflow: hidden;">
      <table v-if="books.length > 0">
        <thead>
          <tr>
            <th width="80">封面</th>
            <th>书名 / ISBN</th>
            <th>作者 / 出版社</th>
            <th>库存状态</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="book in books" :key="book.id">
            <td>
              <img :src="book.cover_url || 'https://via.placeholder.com/150x200?text=No+Cover'" class="book-cover">
            </td>
            <td>
              <div style="font-weight: 600;">{{ book.title }}</div>
              <div style="color: var(--text-secondary); font-size: 0.8rem;">{{ book.isbn }}</div>
            </td>
            <td>
              <div>{{ book.author }}</div>
              <div style="color: var(--text-secondary); font-size: 0.8rem;">{{ book.publisher }}</div>
            </td>
            <td>
              <!-- 业务要求：库存展示 -->
              <span v-if="book.available > 0" class="badge badge-success">
                可借 {{ book.available }} / {{ book.stock }}
              </span>
              <span v-else class="badge badge-danger">暂时缺货</span>
            </td>
            <td class="text-right">
              <div class="flex justify-end gap-2">
                <button class="btn btn-sm" v-if="book.available > 0" @click="handleBorrow(book)">
                  借阅
                </button>
                <button class="btn btn-sm btn-secondary" v-else @click="handleReserve(book.id)">
                  预约
                </button>
                <button v-if="isAdmin" class="btn btn-sm btn-danger" @click="handleDelete(book.id)">
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      
      <!-- 无数据展示 -->
      <div v-else class="empty-state">
        <p>暂无相关图书数据</p>
      </div>

      <!-- 分页 -->
      <div class="pagination flex justify-between" v-if="books.length > 0">
        <button class="btn btn-secondary btn-sm" :disabled="params.page <= 1" @click="changePage(-1)">上一页</button>
        <span style="color: var(--text-secondary);">第 {{ params.page }} / {{ totalPages }} 页</span>
        <button class="btn btn-secondary btn-sm" :disabled="params.page >= totalPages" @click="changePage(1)">下一页</button>
      </div>
    </div>

    <!-- 新增弹窗 -->
    <div v-if="showAddModal" class="modal-overlay" @click.self="showAddModal = false">
      <div class="modal">
        <h2 style="margin-bottom: 24px;">📚 图书入库</h2>
        <div class="form-group"><label>书名</label><input v-model="newBook.title" class="input"></div>
        <div class="form-group"><label>作者</label><input v-model="newBook.author" class="input"></div>
        <div class="form-group"><label>ISBN</label><input v-model="newBook.isbn" class="input" placeholder="978-..."></div>
        <div class="flex gap-4">
          <div class="form-group w-full"><label>出版社</label><input v-model="newBook.publisher" class="input"></div>
          <div class="form-group w-full"><label>库存数量</label><input v-model.number="newBook.stock" type="number" class="input"></div>
        </div>
        <div class="form-group"><label>分类ID</label><input v-model.number="newBook.category_id" type="number" class="input"></div>
        
        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showAddModal = false">取消</button>
          <button class="btn" @click="submitAddBook">确认添加</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue';
import { getBooks, borrowBook, reserveBook, deleteBook, addBook } from '../api';
import { getUser } from '../utils/auth';
import { $message } from '../utils/toast';

const books = ref([]);
const totalPages = ref(1);
const params = reactive({ page: 1, limit: 10, title: '', author: '' });
const showAddModal = ref(false);
const newBook = reactive({ title: '', author: '', isbn: '', category_id: 1, publisher: '', stock: 10 });
const user = getUser();
const isAdmin = computed(() => user?.role === 'admin');

const fetchBooks = async () => {
  try {
    const res = await getBooks(params);
    books.value = res.books || [];
    totalPages.value = res.total_pages || 1;
  } catch(e) { 
    $message.error('数据加载失败'); 
    books.value = []; 
  }
};
const changePage = (delta) => {
  params.page += delta;
  fetchBooks();
};

const handleBorrow = async (book) => {
  // 原生 confirm 还是可以保留的，因为它用来做决定
  if (confirm(`确认借阅 《${book.title}》?`)) {
    try {
      await borrowBook({ book_id: book.id });
      $message.success('借阅成功！请按时归还'); // 替换 alert
      fetchBooks();
    } catch(e) {
      // 错误由 request.js 抛出，或者在这里手动捕获
      // 如果 request.js 里也用了 alert，建议去 request.js 里把 alert 换成 $message.error
    }
  }
};

const handleReserve = async (bookId) => {
  try { 
    await reserveBook({ book_id: bookId }); 
    $message.success('预约成功！书归还后会通知您'); 
  } catch (e) {}
};

const handleDelete = async (id) => {
  if (confirm('确定删除此书?')) {
    await deleteBook(id);
    fetchBooks();
  }
};

const submitAddBook = async () => {
  try {
    await addBook(newBook);
    $message.success('图书入库成功');
    showAddModal.value = false;
    fetchBooks();
  } catch(e) {}
};

onMounted(fetchBooks);
</script>

<style scoped>
.empty-state { padding: 40px; text-align: center; color: var(--text-secondary); }
.pagination { padding: 16px; border-top: 1px solid var(--border); }
</style>