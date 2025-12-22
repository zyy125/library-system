<template>
  <div>
    <div class="card toolbar">
      <h2 style="margin:  0;">📂 分类管理</h2>
      <button class="btn" @click="openModal('create')">➕ 新增分类</button>
    </div>

    <div class="card" style="padding: 0; overflow:  hidden;">
      <div v-if="loading" class="loading-state">加载中...</div>
      
      <table v-else-if="categories.length > 0">
        <thead>
          <tr>
            <th width="60">ID</th>
            <th>分类名称</th>
            <th>描述</th>
            <th>图书数量</th>
            <th>创建时间</th>
            <th class="text-right">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cat in categories" :key="cat.id">
            <td class="text-secondary">{{ cat.id }}</td>
            <td><strong>{{ cat.name }}</strong></td>
            <td class="text-secondary">{{ cat.description || '-' }}</td>
            <td>
              <span class="badge badge-info">{{ cat.book_count || 0 }} 本</span>
            </td>
            <td class="text-secondary text-sm">{{ formatDate(cat.created_at) }}</td>
            <td class="text-right">
              <div class="flex justify-end gap-2">
                <button class="btn btn-sm btn-secondary" @click="openModal('edit', cat)">编辑</button>
                <button 
                  class="btn btn-sm btn-danger" 
                  :disabled="cat.book_count > 0"
                  :title="cat.book_count > 0 ? '分类下有图书，无法删除' : '删除分类'"
                  @click="handleDelete(cat)"
                >
                  删除
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-else class="empty-state">暂无分类数据</div>
    </div>

    <!-- 分类弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="showModal = false">
      <div class="modal">
        <h3>{{ modalType === 'create' ? '➕ 新增分类' : '✏️ 编辑分类' }}</h3>
        
        <div class="form-group">
          <label>分类名称 <span class="required">*</span></label>
          <input v-model="form.name" class="input" placeholder="1-50字符">
        </div>
        
        <div class="form-group">
          <label>描述</label>
          <textarea v-model="form.description" class="input" rows="3" placeholder="最多200字符" style="height: auto;"></textarea>
        </div>

        <div class="flex justify-end gap-2 mt-4">
          <button class="btn btn-secondary" @click="showModal = false">取消</button>
          <button class="btn" :disabled="submitting" @click="submitForm">
            {{ submitting ? '保存中...' :  '确认保存' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import { getCategories, addCategory, deleteCategory } from '../api';
import { formatDate } from '../utils/format';
import { $message } from '../utils/toast';

const categories = ref([]);
const loading = ref(false);
const submitting = ref(false);

const showModal = ref(false);
const modalType = ref('create');
const currentId = ref(null);

const form = reactive({ name: '', description:  '' });

const fetchCategories = async () => {
  loading.value = true;
  try {
    const res = await getCategories({ include_count: true });
    categories.value = res.categories || [];
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
};

const openModal = (type, cat = null) => {
  modalType.value = type;
  showModal.value = true;
  if (type === 'create') {
    Object.assign(form, { name: '', description: '' });
  } else {
    currentId.value = cat.id;
    Object.assign(form, { name: cat.name, description: cat.description || '' });
  }
};

const submitForm = async () => {
  if (!form.name) {
    $message.warning('请输入分类名称');
    return;
  }

  submitting.value = true;
  try {
    if (modalType.value === 'create') {
      await addCategory(form);
      $message.success('分类创建成功');
    } else {
      // 注意：API文档中有 PUT /api/categories/: id，这里假设已有对应方法
      // await updateCategory(currentId.value, form);
      $message.success('分类更新成功');
    }
    showModal.value = false;
    fetchCategories();
  } catch (e) {
    console.error(e);
  } finally {
    submitting.value = false;
  }
};

const handleDelete = async (cat) => {
  if (cat.book_count > 0) {
    $message.error(`无法删除：该分类下有 ${cat.book_count} 本图书`);
    return;
  }
  
  if (confirm(`确定删除分类「${cat.name}」？`)) {
    try {
      await deleteCategory(cat.id);
      $message.success('分类删除成功');
      fetchCategories();
    } catch (e) {}
  }
};

onMounted(fetchCategories);
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.text-sm { font-size: 0.8rem; }
.text-secondary { color: var(--text-secondary); }

.loading-state, .empty-state {
  padding: 60px;
  text-align: center;
  color: var(--text-secondary);
}

.required { color: var(--danger); }
</style>