(function ($) {
    'use strict';
    $(function () {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');
        $('.todo-list-add-btn').on("click", function (event) {
            event.preventDefault();

            var item = $(this).prevAll('.todo-list-input').val();

            if (item) {
                $.post("/todos", {name:item}, addItem)
            
                // todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
                todoListInput.val("");
            }

        });

        var addItem = function(item) {
            if(item.completed) {
                todoListItem.append("<li class='completed'" + " id= '" + item.id + "'>" + "<div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked'/>" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            } else {
                todoListItem.append("<li " + " id= '" + item.id + "'>" + "<div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
                 
        }

        $.get('/todos', function(items) {
            items.forEach(e => {
                addItem(e)
            });
        });

        todoListItem.on('change', '.checkbox', function () {
            var id = $(this).closest('li').attr('id');
            var $self = $(this)
            var completed = true
            if ($self.attr('checked')) {
                completed = false
            } 
            $.get("complete-todo/" + id + "?complete=" + completed, function(data) {
                
                if (completed) {
                    $self.attr('checked', 'checked');
                } else {
                    $self.removeAttr('checked');
                }
    
                $self.closest("li").toggleClass('completed');
         
            }) 

           
        });

        todoListItem.on('click', '.remove', function () {
            var id = $(this).closest('li').attr('id');
            // func 이 불릴 때의 this가 될 수 있으므로 미리 저장
            var $self = $(this)
            //서버 요청은 ajax *get, post 제외
            $.ajax({
                url: "todos/" + id,
                type: "DELETE",
                success: function(data) {
                    if(data.success) {
                        $self.parent().remove()
                    }
                }
            })
            // $(this).parent().remove();
        });

    });
})(jQuery);